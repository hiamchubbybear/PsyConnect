import asyncio
import random
import string
import csv
import json
import os
from datetime import datetime
from playwright.async_api import async_playwright

BASE_URL = "http://localhost:4200"
LOG_FILE_CSV = "automation/provisioned_users.csv"
LOG_FILE_JSON = "automation/provisioned_users.json"

def log_user_csv(user):
    file_exists = os.path.isfile(LOG_FILE_CSV)
    os.makedirs(os.path.dirname(LOG_FILE_CSV), exist_ok=True)
    with open(LOG_FILE_CSV, mode='a', newline='', encoding='utf-8') as f:
        writer = csv.DictWriter(f, fieldnames=["timestamp", "role", "username", "password", "first_name", "last_name"])
        if not file_exists:
            writer.writeheader()

        writer.writerow({
            "timestamp": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            "role": user["role"],
            "username": user["username"],
            "password": user["password"],
            "first_name": user["first_name"],
            "last_name": user["last_name"]
        })
    print(f"📖 Account logged to CSV: {LOG_FILE_CSV}")

def log_user_json(user):
    os.makedirs(os.path.dirname(LOG_FILE_JSON), exist_ok=True)
    users = []
    if os.path.isfile(LOG_FILE_JSON):
        try:
            with open(LOG_FILE_JSON, 'r', encoding='utf-8') as f:
                users = json.load(f)
        except:
            users = []

    users.append({
        "timestamp": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        "role": user["role"],
        "username": user["username"],
        "password": user["password"],
        "first_name": user["first_name"],
        "last_name": user["last_name"]
    })

    with open(LOG_FILE_JSON, 'w', encoding='utf-8') as f:
        json.dump(users, f, indent=2, ensure_ascii=False)
    print(f"📦 Account logged to JSON: {LOG_FILE_JSON}")

# Lists for randomization
FIRST_NAMES = ["An", "Binh", "Chi", "Dung", "Em", "Giang", "Hung", "Khoa", "Linh", "Minh", "Nam", "Phuong", "Quang", "Son", "Thao", "Uyen", "Viet", "Xuan", "Yen", "Duy"]
LAST_NAMES = ["Nguyen", "Tran", "Le", "Pham", "Hoang", "Phan", "Vu", "Dang", "Bui", "Do", "Ho", "Ngo", "Duong", "Ly"]
CITIES = ["Ho Chi Minh City", "Ha Noi", "Da Nang", "Can Tho", "Hai Phong", "Hue", "Nha Trang", "Da Lat"]
UNIVERSITIES = ["University of Medicine", "National University", "International University", "FPT University", "RMIT University", "Ton Duc Thang University", "Hoa Sen University"]
FIELDS = ["Psychology", "Counseling", "Mental Health", "Clinical Psychology", "Behavioral Therapy", "Social Work"]

def generate_random_string(length=8):
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=length))

async def create_user(page, role="Client", index=1):
    first_name = random.choice(FIRST_NAMES)
    last_name = random.choice(LAST_NAMES)
    username = f"{last_name.lower()}_{first_name.lower()}_{generate_random_string(4)}_{index}"
    email = f"{username}@example.com"
    password = "Password123!"

    # Random DOB between 18 and 60 years ago
    year = random.randint(1965, 2005)
    month = random.randint(1, 12)
    day = random.randint(1, 28)
    dob = f"{year}-{month:02d}-{day:02d}"

    gender = random.choice(["Male", "Female", "None"])
    city = random.choice(CITIES)

    print(f"🚀 Creating {role} account: {username} ({first_name} {last_name})")

    await page.goto(f"{BASE_URL}/auth/signup")

    # Step 1: Personal Info
    await page.fill('app-input[formControlName="firstName"] input', first_name)
    await page.fill('app-input[formControlName="lastName"] input', last_name)

    # Step 2: DOB & Gender
    await page.fill('app-input[formControlName="dateOfBirth"] input', dob)
    await page.select_option('select[formControlName="gender"]', value=gender)

    # Step 3: Address
    await page.fill('app-input[formControlName="address"] input', city)
    await page.keyboard.press("Tab")

    # Step 4: Email
    await page.fill('app-input[formControlName="email"] input', email)

    # Step 5: Role
    role_selector = 'input[value="Therapist"]' if role == "Therapist" else 'input[value="Client"]'
    await page.dispatch_event(role_selector, 'click')

    # Step 6: Account Details
    await page.fill('app-input[formControlName="username"] input', username)
    await page.fill('app-input[formControlName="password"] input', password)
    await page.fill('app-input[formControlName="confirmPassword"] input', password)
    await page.dispatch_event('app-checkbox[formControlName="acceptTerms"] input', 'click')

    # Submit
    await page.click('button[type="submit"]')

    # Wait for redirect (either login or activate)
    await page.wait_for_function("() => window.location.href.includes('/auth/login') || window.location.href.includes('/activate')", timeout=10000)

    print(f"✅ User {username} signed up!")
    return {"username": username, "password": password, "role": role, "first_name": first_name, "last_name": last_name}

async def login(page, user):
    print(f"🔑 Logging in as {user['username']}...")
    await page.goto(f"{BASE_URL}/auth/login")

    if await page.is_visible('.btn-circle-mail'):
        await page.click('.btn-circle-mail')

    await page.fill('app-input[formControlName="username"] input', user['username'])
    await page.fill('app-input[formControlName="password"] input', user['password'])
    await page.click('button[type="submit"]')

    await page.wait_for_url(f"{BASE_URL}/feature/feed", timeout=10000)
    print(f"✅ Logged in successfully!")

async def setup_consultation_profile(page, user):
    print(f"🛠️ Setting up consultation profile for {user['username']}...")
    await page.goto(f"{BASE_URL}/feature/consultation/create-profile")

    # Wait for Step 1
    await page.wait_for_selector('.active-pane h4:has-text("Consultation Modes")', timeout=10000)

    # Step 1: Configuration
    cards = page.locator('.active-pane .selectable-card')
    count = await cards.count()
    if count > 0:
        mode_index = random.randint(0, count - 1)
        await cards.nth(mode_index).click()

    price = random.randint(200, 2000) * 1000
    await page.fill('.active-pane input[formControlName="rage_price"]', str(price))

    if user['role'].lower() == "therapist":
        await page.select_option('.active-pane select[formControlName="currency"]', "VND")
    await page.click('.btn-next')

    # Step 2: Attributes
    await page.wait_for_selector('.active-pane h4:has-text("Languages")', timeout=10000)

    lang_cards = page.locator('.active-pane .section-group:first-of-type .selectable-card')
    lang_count = await lang_cards.count()
    if lang_count > 0:
        await lang_cards.nth(random.randint(0, lang_count - 1)).click()

    topic_cards = page.locator('.active-pane .section-group:nth-of-type(2) .selectable-card')
    topic_count = await topic_cards.count()
    if topic_count > 0:
        await topic_cards.nth(random.randint(0, topic_count - 1)).click()

    await page.click('.btn-next')

    # Step 3: Details
    await page.wait_for_selector('.active-pane input[formControlName="address"]', timeout=10000)

    if user['role'].lower() == "therapist":
        await page.fill('.active-pane input[formControlName="name"]', f"{user['last_name']} {user['first_name']}")

    city = random.choice(CITIES)
    await page.fill('.active-pane input[formControlName="address"]', city)

    if user['role'].lower() == "therapist":
        await page.select_option('.active-pane select[formControlName="professional_title_code"]', index=random.randint(1, 4))
        await page.fill('.active-pane input[formControlName="experience"]', str(random.randint(1, 20)))
        await page.click('.btn-next')

        # Step 4: Professional Info
        await page.wait_for_selector('.active-pane h4:has-text("Degrees")', timeout=10000)

        num_degrees = random.randint(1, 2)
        for d in range(num_degrees):
            await page.click('.active-pane .btn-add >> nth=0')
            await page.select_option(f'.active-pane select[formControlName="type"] >> nth={d}', index=random.randint(1, 4))
            await page.fill(f'.active-pane input[formControlName="field"] >> nth={d}', random.choice(FIELDS))
            await page.fill(f'.active-pane input[formControlName="institution"] >> nth={d}', random.choice(UNIVERSITIES))
            await page.fill(f'.active-pane input[formControlName="year"] >> nth={d}', str(random.randint(2010, 2023)))

        await page.click('.active-pane .btn-add >> nth=1')
        await page.fill('.active-pane [formArrayName="certifications"] input[formControlName="name"] >> nth=0', "Mental Health Professional Certificate")
        await page.fill('.active-pane [formArrayName="certifications"] input[formControlName="issuer"] >> nth=0', "Global Health Org")
        await page.fill('.active-pane [formArrayName="certifications"] input[formControlName="year"] >> nth=0', "2021")

    # Submit Profile
    submit_btn = page.locator('.header-actions button:has-text("Save"), .header-actions button:has-text("Update"), .header-actions button:has-text("Profile")')
    await submit_btn.click()

    # FIX: The app redirects to the feed/homepage after saving profile
    print("⏳ Waiting for redirect after saving profile...")
    await page.wait_for_function("() => window.location.href.includes('/feed') || window.location.pathname === '/'", timeout=15000)
    print(f"✨ Consultation profile for {user['username']} is READY!")

async def clear_app_state(page, context):
    """Thoroughly clear all application state to ensure logout."""
    print("🧹 Clearing application state (cookies, localStorage, sessionStorage)...")
    await context.clear_cookies()
    await page.evaluate("window.localStorage.clear()")
    await page.evaluate("window.sessionStorage.clear()")
    # Go to a neutral page to ensure no background redirects catch us
    await page.goto("about:blank")

async def main(num_users=None):
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=False)
        context = await browser.new_context()
        page = await context.new_page()

        success_count = 0
        i = 1
        try:
            while True:
                try:
                    role = random.choice(["Client", "Therapist"])
                    user = await create_user(page, role=role, index=i)

                    # Login and Setup Profile
                    await login(page, user)
                    await setup_consultation_profile(page, user)

                    # Log account details in both formats
                    log_user_csv(user)
                    log_user_json(user)

                    print(f"✅ Successfully finished provision for user {i} ({role})")
                    success_count += 1
                except Exception as e:
                    print(f"\n❌ Error during cycle for user {i}: {e}")
                    try:
                        await page.screenshot(path=f"automation/error_step_{i}.png")
                        print(f"📸 Screenshot saved to automation/error_step_{i}.png")
                    except: pass

                # ALWAYS clear state before next user, even if failed
                await clear_app_state(page, context)
                print(f"--- Cycle {i} completed ---")

                if num_users and success_count >= num_users:
                    break

                i += 1
                delay = random.randint(2, 5)
                print(f"⏳ Waiting {delay}s for next cycle... (Created: {success_count})\n")
                await asyncio.sleep(delay)

        except KeyboardInterrupt:
            print("\n🛑 Automation stopped by user.")
        finally:
            await browser.close()
            print(f"\n🎉 Automation session ended. Total users created: {success_count}")

if __name__ == "__main__":
    import sys
    count = int(sys.argv[1]) if len(sys.argv) > 1 else None
    asyncio.run(main(count))
