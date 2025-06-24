import random
import json

addresses = ["Hanoi", "HCMC", "Danang", "Hue"]
languages_pool = ["Vietnamese", "English", "Japanese", "French"]
specializations_pool = ["stress", "anxiety", "depression", "marriage", "career"]
consultation_modes_pool = ["online", "offline"]
availability_pool = [
    "Morning_Weekday", "Afternoon_Weekday",
    "Evening_Weekday", "Morning_Weekend", "Evening_Weekend"
]

def generate_random_therapist(i):
    return {
        "profile_id": f"T{i:04d}",
        "address": random.choice(addresses),
        "languages": random.sample(languages_pool, k=random.randint(1, 2)),
        "specialization": random.sample(specializations_pool, k=random.randint(1, 3)),
        "consultation_modes": random.sample(consultation_modes_pool, k=1),
        "experience": random.randint(1, 15),
        "rating": round(random.uniform(3.5, 5.0), 1),
        "currency": "VND",
        "range_price": random.randint(300000, 700000),
        "availability": random.sample(availability_pool, k=random.randint(1, 3)),
        "is_available": random.choice([True, True, False]),
        "matched_clients": []
    }

def generate_client():
    return {
        "profile_id": "C999",
        "address": "Hanoi",
        "languages": ["Vietnamese", "English"],
        "issue_detail": "stress",
        "consultation_modes": ["online"],
        "range_price": 500000,
        "availability": ["Morning_Weekday", "Afternoon_Weekday"],
        "preferred_therapist_gender": None,
        "experience_level": "mid",
        "specialization": ["stress", "anxiety"],
        "urgency_level": "medium",
        "session_duration": 60,
        "preferred_therapist_language": "Vietnamese",
        "is_flexible_with_schedule": True
    }
hard_client = {
    "profile_id": "C_HARD",
    "address": "Hue",  
    "languages": ["Japanese"],
    "issue_detail": "career crisis",
    "consultation_modes": ["offline"],
    "range_price": 300000,
    "availability": ["Evening_Weekend"],
    "preferred_therapist_gender": "female",
    "experience_level": "senior",
    "specialization": ["career", "marriage"],
    "urgency_level": "high",
    "session_duration": 90,
    "preferred_therapist_language": "Japanese",
    "is_flexible_with_schedule": False
}

with open("hard_client.json", "w") as f:
    json.dump(hard_client, f, indent=2, ensure_ascii=False)

print("✅ Đã tạo client C_HARD có yêu cầu rất khó match.")

if __name__ == "__main__":
    n = 1000
    therapists = [generate_random_therapist(i) for i in range(n)]
    client = generate_client()


    with open("generated_therapists.json", "w") as f:
        json.dump(therapists, f, indent=2, ensure_ascii=False)
    with open("generated_client.json", "w") as f:
        json.dump(client, f, indent=2, ensure_ascii=False)

    print(f"✅ Đã tạo {n} therapists và 1 client để test.")
