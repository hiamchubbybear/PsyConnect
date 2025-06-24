# Debug mode
# Therapist matching using traditional ML (similarity-based scoring)

import pandas as pd
import numpy as np
from sklearn.preprocessing import MultiLabelBinarizer
from sklearn.metrics.pairwise import cosine_similarity
import json
import time

start_time = time.time()

# client = {
#     "profile_id": "C002",
#     "address": "HCMC",
#     "languages": ["English"],
#     "issue_detail": "marriage",
#     "consultation_modes": ["online"],
#     "range_price": 600000,
#     "availability": ["Morning_Weekday"],
#     "preferred_therapist_gender": None,
#     "experience_level": "senior",
#     "specialization": ["marriage"],
#     "urgency_level": "low",
#     "session_duration": 60,
#     "preferred_therapist_language": "English",
#     "is_flexible_with_schedule": False
# }


# therapists = [
#     {
#         "profile_id": "T001",
#         "address": "Hanoi",
#         "languages": ["Vietnamese"],
#         "specialization": ["stress", "depression"],
#         "consultation_modes": ["online"],
#         "experience": 5,
#         "rating": 4.7,
#         "currency": "VND",
#         "rage_price": 400000,
#         "availability": ["Morning_Weekday", "Afternoon_Weekday"],
#         "is_available": True,
#         "matched_clients": []
#     },
#     {
#         "profile_id": "T002",
#         "address": "HCMC",
#         "languages": ["English"],
#         "specialization": ["marriage"],
#         "consultation_modes": ["online"],
#         "experience": 10,
#         "rating": 4.5,
#         "currency": "VND",
#         "rage_price": 600000,
#         "availability": ["Morning_Weekday", "Afternoon_Weekday"],
#         "is_available": True,
#         "matched_clients": []
#     },
#     {
#         "profile_id": "T003",
#         "address": "Hanoi",
#         "languages": ["Vietnamese", "English"],
#         "specialization": ["stress", "anxiety"],
#         "consultation_modes": ["online"],
#         "experience": 3,
#         "rating": 4.9,
#         "currency": "VND",
#         "rage_price": 450000,
#         "availability": ["Morning_Weekday", "Afternoon_Weekday"],
#         "is_available": True,
#         "matched_clients": []
#     }
# ]
with open("generated_therapists.json") as f:
    therapists = json.load(f)
with open("hard_client.json") as f:
    client = json.load(f)


def extract_features(client, therapists):
    df = pd.DataFrame(therapists)
    print("\nInitial DataFrame:")
    print(df)

    mlb = MultiLabelBinarizer()

    for field in ["languages", "specialization", "consultation_modes", "availability"]:
        df_field = mlb.fit_transform(df[field])
        print(f"\nEncoded therapist {field}:")
        print(pd.DataFrame(df_field, columns=[f"{field}_{cls}" for cls in mlb.classes_]))

        df = df.join(pd.DataFrame(df_field, columns=[f"{field}_{cls}" for cls in mlb.classes_]))

        client_encoded = mlb.transform([client[field]])
        print(f"Encoded client {field}:")
        print(client_encoded)

        df[f"{field}_sim"] = cosine_similarity(df_field, client_encoded).flatten()
        print(f"{field}_sim:")
        print(df[f"{field}_sim"])

    df["address_sim"] = df["address"].apply(lambda x: 1.0 if x == client["address"] else 0.5)
    print("\nAddress similarity:")
    print(df["address_sim"])

    df["price_match"] = df["range_price"].apply(lambda x: 1.0 if x <= client["range_price"] else 0.0)
    print("\nPrice match:")
    print(df["price_match"])

    df["exp_score"] = df["experience"] / df["experience"].max()
    print("\nExperience score:")
    print(df["exp_score"])

    return df

def compute_match_score(df):
    weights = {
        "languages_sim": 0.2,
        "specialization_sim": 0.3,
        "consultation_modes_sim": 0.1,
        "availability_sim": 0.15,
        "address_sim": 0.1,
        "price_match": 0.1,
        "exp_score": 0.05
    }

    df["match_score"] = sum(df[k] * w for k, w in weights.items()) * 100
    print("\nFinal match scores:")
    print(df[["profile_id"] + list(weights.keys()) + ["match_score"]])
    return df

df_feat = extract_features(client, therapists)
df_scored = compute_match_score(df_feat)
df_scored = df_scored[df_scored["is_available"] == True].sort_values("match_score", ascending=False)
print("\nSorted therapists by match_score:")
print(df_scored[["profile_id", "match_score"]])
# Ghi tất cả therapist có match_score vào file JSON
detailed_output = []
for _, row in df_scored.iterrows():
    detailed_output.append({
        "profile_id": row["profile_id"],
        "match_score": round(row["match_score"], 2),
        "languages_sim": round(row["languages_sim"], 2),
        "specialization_sim": round(row["specialization_sim"], 2),
        "availability_sim": round(row["availability_sim"], 2),
        "consultation_modes_sim": round(row["consultation_modes_sim"], 2),
        "address_sim": round(row["address_sim"], 2),
        "price_match": row["price_match"],
        "exp_score": round(row["exp_score"], 2),
        "is_available": row["is_available"]
    })

with open("matching_results.json", "w", encoding="utf-8") as f:
    json.dump(detailed_output, f, indent=2, ensure_ascii=False)

print("✅ Đã xuất kết quả matching vào matching_results.json")

output = []
for _, row in df_scored.head(3).iterrows():
    output.append({
        "profile_id": row["profile_id"],
        "match_score": round(row["match_score"], 2),
        "reason": "Matched on specialization, language, availability"
    })
print("\nFinal Output:")
print(json.dumps(output, indent=2, ensure_ascii=False))
end_time = time.time()
print(f"\n⏱️ Tổng thời gian thực thi: {end_time - start_time:.4f} giây")
