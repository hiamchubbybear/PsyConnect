import asyncio
import random
import json
import os
from datetime import datetime
from playwright.async_api import async_playwright

# Production URL targeting
BASE_URL = "https://psyconnect.chessy.dev"
USER_FILE = "automation/automation/provisioned_users.json"

# Fallback search if the above nested directory doesn't exist
if not os.path.exists(USER_FILE):
    USER_FILE = "automation/provisioned_users.json"

# Psychology keywords for Wikipedia fetch
WIKI_KEYWORDS = [
    "Cognitive behavioral therapy", "Attachment theory", "Neuroplasticity",
    "Social psychology", "Defense mechanism", "Unconscious mind",
    "Positive psychology", "Developmental psychology", "Clinical psychology",
    "Behavioral therapy", "Emotional intelligence", "Maslow's hierarchy of needs",
    "Erikson's stages of psychosocial development", "Classical conditioning",
    "Operant conditioning", "Gestalt psychology", "Humanistic psychology",
    "Psychology", "Mental health", "Psychotherapy"
]

async def fetch_wikipedia_content(context, keyword):
    """Fetch summary from Wikipedia API using Playwright's built-in request API."""
    url = f"https://en.wikipedia.org/api/rest_v1/page/summary/{keyword.replace(' ', '_')}"
    print(f"📡 Fetching psychology content for: '{keyword}'...")
    try:
        response = await context.request.get(url)
        if response.status == 200:
            data = await response.json()
            return {
                "title": data.get("title", keyword),
                "description": data.get("extract", "No summary available.")
            }
        else:
            print(f"⚠️ Wikipedia API returned status {response.status}")
    except Exception as e:
        print(f"❌ Error fetching from Wikipedia: {e}")
    return None

async def clear_app_state(page, context):
    """Thoroughly clear all application state to ensure logout."""
    print("🧹 Clearing application state...")
    await context.clear_cookies()
    await page.evaluate("window.localStorage.clear()")
    await page.evaluate("window.sessionStorage.clear()")
    await page.goto("about:blank")

async def login(page, user):
    print(f"🔑 Logging in to production as {user['username']}...")
    await page.goto(f"{BASE_URL}/auth/login")

    # Handle optional circle mail button if present
    try:
        if await page.is_visible('.btn-circle-mail', timeout=2000):
            await page.click('.btn-circle-mail')
    except: pass

    await page.fill('app-input[formControlName="username"] input', user['username'])
    await page.fill('app-input[formControlName="password"] input', user['password'])
    await page.click('button[type="submit"]')

    # Wait for feed
    await page.wait_for_url(f"{BASE_URL}/feature/feed", timeout=15000)
    print(f"✅ Logged in successfully!")

async def create_post(page, title, description, anonymous=False):
    print(f"📝 Publishing article on production: '{title}' (Anon: {anonymous})")
    await page.goto(f"{BASE_URL}/feature/article/create")

    # Wait for the form to be ready
    await page.wait_for_selector('input#title', timeout=15000)

    # Fill Title
    await page.fill('input#title', title)

    # Fill Content (ID is 'content' in production template)
    await page.fill('textarea#content', description)

    # Toggle Anonymous if requested AND if the element exists in this version
    if anonymous:
        try:
            # Check if anon checkbox exists (it might be in a different template version)
            anon_selector = 'input[name="anonymous"]'
            if await page.is_visible(anon_selector, timeout=2000):
                print("👤 Enabling anonymous mode...")
                await page.dispatch_event(anon_selector, 'click')
        except: pass

    # Submit
    await page.click('button[type="submit"]')

    # Wait for completion (production redirects to feed)
    print("⏳ Finalizing post publication...")
    await page.wait_for_url("**/feature/feed", timeout=20000)
    print(f"✨ Article version of '{title}' is LIVE on production!")

async def main(num_posts=None):
    if not os.path.exists(USER_FILE):
        print(f"❌ Error: {USER_FILE} not found. Run the provision_users.py script first.")
        return

    async with async_playwright() as p:
        # Launch browser
        browser = await p.chromium.launch(headless=False)
        context = await browser.new_context()
        page = await context.new_page()

        try:
            with open(USER_FILE, 'r', encoding='utf-8') as f:
                users = json.load(f)
        except Exception as e:
            print(f"❌ Error reading users: {e}")
            return

        if not users:
            print("❌ No users found in provisioned_users.json.")
            return

        posts_created = 0
        cycle = 1

        try:
            while True:
                try:
                    user = random.choice(users)
                    keyword = random.choice(WIKI_KEYWORDS)

                    # Fetch content using context.request
                    content = await fetch_wikipedia_content(context, keyword)

                    if not content:
                        print("⏭️ Skipping cycle due to fetch failure.")
                        continue

                    is_anon = random.choice([True, False, False]) # 33% chance anon

                    await login(page, user)
                    await create_post(page, content['title'], content['description'], anonymous=is_anon)

                    posts_created += 1
                    print(f"🚀 Finished production cycle {cycle} by {user['username']}")
                except Exception as e:
                    print(f"\n❌ Production Error during cycle {cycle}: {e}")
                    try:
                        os.makedirs("automation/screenshots", exist_ok=True)
                        await page.screenshot(path=f"automation/screenshots/prod_error_{cycle}.png")
                        print(f"📸 Debug screenshot: automation/screenshots/prod_error_{cycle}.png")
                    except: pass

                await clear_app_state(page, context)

                if num_posts and posts_created >= num_posts:
                    break

                cycle += 1
                delay = random.randint(10, 30) # Slower delay for production
                print(f"⏳ Waiting {delay}s for next worker cycle... (Production Posts: {posts_created})\n")
                await asyncio.sleep(delay)

        except KeyboardInterrupt:
            print("\n🛑 Production automation stopped by user.")
        finally:
            await browser.close()
            print(f"\n🎉 Production session ended. Total posts created: {posts_created}")

if __name__ == "__main__":
    import sys
    count = int(sys.argv[1]) if len(sys.argv) > 1 else None
    asyncio.run(main(count))
