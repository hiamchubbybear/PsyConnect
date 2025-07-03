from flask import Flask, request, jsonify
from sklearn.preprocessing import MultiLabelBinarizer
from sklearn.metrics.pairwise import cosine_similarity
import pandas as pd

app = Flask(__name__)

@app.route('/recommend', methods=['POST'])
def recommend():
    try:
        print("Received request")
        data = request.json
        print("Raw data:", data)

        client = data.get("clientRaw")
        therapists = data.get("therapistsRaw")

        if not client or not therapists:
            return jsonify({"error": "Missing client or therapist data"}), 400

        df = pd.DataFrame(therapists)

        for field in ["languages", "specialization", "consultation_modes"]:
            if field not in df.columns:
                print(f"⚠️ Therapist data missing field: {field}. Skipping...")
                continue

            # Đảm bảo field là list
            df[field] = df[field].apply(lambda x: x if isinstance(x, list) else [])

            try:
                mlb = MultiLabelBinarizer()
                df_field = mlb.fit_transform(df[field])
                df = df.join(pd.DataFrame(df_field, columns=[f"{field}_{cls}" for cls in mlb.classes_]))
                client_filtered = [val for val in client.get(field, []) if val in mlb.classes_]
                client_encoded = mlb.transform([client_filtered])
                df[f"{field}_sim"] = cosine_similarity(df_field, client_encoded).flatten()
            except Exception as e:
                print(f"❌ Error processing field {field}: {e}")
                df[f"{field}_sim"] = 0.0
        # Availability similarity
        df["availability"] = df["availability"].apply(lambda x: x["days"] + x["time_slots"])
        client_avail = client["availability"]["days"] + client["availability"]["time_slots"]
        mlb = MultiLabelBinarizer()
        df_field = mlb.fit_transform(df["availability"])
        df = df.join(pd.DataFrame(df_field, columns=[f"availability_{cls}" for cls in mlb.classes_]))
        client_avail_filtered = [val for val in client_avail if val in mlb.classes_]
        client_encoded = mlb.transform([client_avail_filtered])
        df["availability_sim"] = cosine_similarity(df_field, client_encoded).flatten()

        # Address similarity
        df["address_sim"] = df["address"].apply(lambda x: 1.0 if x == client["address"] else 0.5)

        # Price match (tạm set mặc định vì không có trường giá)
        df["price_match"] = 1.0

        # Experience score
        df["exp_score"] = df["experience"] / df["experience"].max()

        # Nếu thiếu cột is_available thì giả sử tất cả đều available
        if "is_available" not in df.columns:
            df["is_available"] = True

        # Match score với trọng số
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
        df = df[df["is_available"] == True].sort_values("match_score", ascending=False)

        result = {
            "client_id": client["profile_id"],
            "swipes": []
        }

        for _, row in df.iterrows():
            reasons = []
            if row["specialization_sim"] > 0: reasons.append("Matched specialization")
            if row["languages_sim"] > 0: reasons.append("Matched language")
            if row["availability_sim"] > 0: reasons.append("Matched availability")

            result["swipes"].append({
                "therapist_id": row["profile_id"],
                "points": round(row["match_score"], 2),
                "reasons": reasons
            })
        print(result)
        return jsonify(result)

    except Exception as e:
        print("Error occurred:", e)
        return jsonify({"error": "Internal server error", "details": str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True)
