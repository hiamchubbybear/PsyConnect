from flask import Flask, request, jsonify
from sklearn.preprocessing import MultiLabelBinarizer
from sklearn.metrics.pairwise import cosine_similarity
import pandas as pd
from datetime import datetime

app = Flask(__name__)

@app.route('/recommend', methods=['POST'])
def recommend():
    try:
        print("✅ Received request")
        data = request.json
        print("📦 Raw data:", data)

        client = data.get("clientRaw")
        therapists = data.get("therapistsRaw")

        if not client or not therapists:
            return jsonify({"error": "Missing client or therapist data"}), 400

        df = pd.DataFrame(therapists)


        weights = {
            "languages_sim": 0.2,
            "specialization_sim": 0.3,
            "consultation_modes_sim": 0.1,
            "availability_sim": 0.15,
            "address_sim": 0.1,
            "price_match": 0.1,
            "exp_score": 0.05
        }


        for field in ["languages", "specialization", "consultation_modes"]:
            if field not in df.columns:
                df[field] = [[] for _ in range(len(df))]

            df[field] = df[field].apply(lambda x: x if isinstance(x, list) else [])

            try:
                mlb = MultiLabelBinarizer()
                df_field = mlb.fit_transform(df[field])
                df = df.join(pd.DataFrame(df_field, columns=[f"{field}_{cls}" for cls in mlb.classes_]))
                client_vals = [val for val in client.get(field, []) if val in mlb.classes_]
                client_encoded = mlb.transform([client_vals])
                df[f"{field}_sim"] = cosine_similarity(df_field, client_encoded).flatten()
            except Exception as e:
                print(f" Error encoding {field}: {e}")
                df[f"{field}_sim"] = 0.0


        df["availability"] = df["availability"].apply(lambda x: x["days"] + x["time_slots"])
        client_avail = client["availability"]["days"] + client["availability"]["time_slots"]
        mlb = MultiLabelBinarizer()
        df_field = mlb.fit_transform(df["availability"])
        df = df.join(pd.DataFrame(df_field, columns=[f"availability_{cls}" for cls in mlb.classes_]))
        client_avail_filtered = [val for val in client_avail if val in mlb.classes_]
        client_encoded = mlb.transform([client_avail_filtered])
        df["availability_sim"] = cosine_similarity(df_field, client_encoded).flatten()


        df["address_sim"] = df["address"].apply(lambda x: 1.0 if x == client["address"] else 0.5)


        df["price_match"] = 1.0


        df["experience"] = df["experience"].fillna(0)
        max_exp = df["experience"].max() if df["experience"].max() > 0 else 1
        df["exp_score"] = df["experience"] / max_exp

        df["is_available"] = df.get("is_available", False)
        df["is_available"] = df["is_available"].fillna(False)

        swipes = []
        now = datetime.utcnow().isoformat() + "Z"


        for _, row in df.iterrows():
            therapist_id = row["profile_id"]
            if not row.get("is_available", True):
                swipes.append({
                    "client_id": client["profile_id"],
                    "therapist_id": therapist_id,
                    "points": 0.0,
                    "reasons": ["Therapist not available"],
                    "status": "pending",
                    "created_at": now
                })
                continue

            score = sum(row.get(k, 0.0) * w for k, w in weights.items()) * 100

            reasons = []
            if row.get("specialization_sim", 0) > 0:
                reasons.append("Matched specialization")
            if row.get("languages_sim", 0) > 0:
                reasons.append("Matched language")
            if row.get("availability_sim", 0) > 0:
                reasons.append("Matched availability")

            swipes.append({
                "client_id": client["profile_id"],
                "therapist_id": therapist_id,
                "points": round(score, 2),
                "reasons": reasons,
                "status": "pending",
                "created_at": now
            })

        # Sắp xếp theo điểm cao nhất
        swipes.sort(key=lambda x: x["points"], reverse=True)

        result = {
            "client_id": client["profile_id"],  # 🔧 Thêm dòng này để Go decode được
            "swipes": swipes
        }
        print(" Response:", result)
        return jsonify(result)

    except Exception as e:
        print(" Internal error:", e)
        return jsonify({"error": "Internal server error", "details": str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True)
