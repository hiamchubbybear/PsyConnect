import pandas as pd
import numpy as np
import difflib
from sklearn.preprocessing import MultiLabelBinarizer
from sklearn.metrics.pairwise import cosine_similarity
from datetime import datetime

def calculate_similarity(client, therapists):
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

    # Binarizer fields
    for field in ["languages", "specialization", "consultation_modes"]:
        if field not in df.columns:
            df[field] = [[] for _ in range(len(df))]

        df[field] = df[field].apply(lambda x: x if isinstance(x, list) else [])

        try:
            mlb = MultiLabelBinarizer()
            df_field = mlb.fit_transform(df[field])
            client_vals = [val for val in client.get(field, []) if val in mlb.classes_]
            if not client_vals:
                df[f"{field}_sim"] = 0.0
            else:
                client_encoded = mlb.transform([client_vals])
                df[f"{field}_sim"] = cosine_similarity(df_field, client_encoded).flatten()
        except Exception as e:
            df[f"{field}_sim"] = 0.0

    # Availability
    if "availability" not in df.columns:
        df["availability"] = [{} for _ in range(len(df))]
        
    df["availability_combined"] = df["availability"].apply(lambda x: x.get("days", []) + x.get("time_slots", []) if isinstance(x, dict) else [])
    
    client_avail_dict = client.get("availability", {})
    if not isinstance(client_avail_dict, dict):
        client_avail_dict = {}
    client_avail = client_avail_dict.get("days", []) + client_avail_dict.get("time_slots", [])
    
    try:
        mlb = MultiLabelBinarizer()
        df_field = mlb.fit_transform(df["availability_combined"])
        client_avail_filtered = [val for val in client_avail if val in mlb.classes_]
        if not client_avail_filtered:
            df["availability_sim"] = 0.0
        else:
            client_encoded = mlb.transform([client_avail_filtered])
            df["availability_sim"] = cosine_similarity(df_field, client_encoded).flatten()
    except Exception as e:
        df["availability_sim"] = 0.0

    # Address Sim (Fuzzy)
    client_address = str(client.get("address", "")).lower()
    df["address"] = df.get("address", "")
    
    def address_similarity(therapist_address):
        t_addr = str(therapist_address).lower()
        if not t_addr or not client_address:
            return 0.5
        if t_addr in client_address or client_address in t_addr:
            return 1.0
        return difflib.SequenceMatcher(None, client_address, t_addr).ratio()
    
    df["address_sim"] = df["address"].apply(address_similarity)

    # Price Match
    client_price = client.get("rage_price", 0) # JSON field matches Go's 'rage_price' typo
    df["rage_price"] = df.get("rage_price", 0)
    
    def price_similarity(therapist_price):
        if pd.isna(therapist_price) or therapist_price <= 0:
            return 1.0 # assume affordable if unknown
        if client_price <= 0:
            return 0.5 # unknown client budget
        if therapist_price <= client_price:
            return 1.0
        # Decay: if price is 2x budget, sim = 0.0
        ratio = therapist_price / client_price
        return max(0.0, 1.0 - (ratio - 1.0))
    
    df["price_match"] = df["rage_price"].apply(price_similarity)

    # Experience Score (Log scaled)
    df["experience"] = df.get("experience", 0).fillna(0)
    max_exp = df["experience"].max()
    if max_exp > 0:
        df["exp_score"] = np.log1p(df["experience"]) / np.log1p(max_exp)
    else:
        df["exp_score"] = 0.0

    df["is_available"] = df.get("is_available", False)
    df["is_available"] = df["is_available"].fillna(False)

    swipes = []
    now = datetime.utcnow().isoformat() + "Z"

    for _, row in df.iterrows():
        therapist_id = row.get("profile_id")
        if not therapist_id or pd.isna(therapist_id):
            continue
            
        if not row.get("is_available", True):
            continue # Removed pending 0.0 generation so db isn't flooded with unavailable therapists.

        score = sum(row.get(k, 0.0) * w for k, w in weights.items()) * 100

        reasons = []
        if row.get("specialization_sim", 0) > 0.5:
            reasons.append("Matched specialization")
        if row.get("languages_sim", 0) > 0.5:
            reasons.append("Matched language")
        if row.get("availability_sim", 0) > 0.5:
            reasons.append("Matched availability")
        if row.get("price_match", 0) >= 0.8:
            reasons.append("Within your budget")

        penalty = 1.0
        # Apply penalty if missing data
        if row.get("specialization_sim", 0) == 0 and row.get("languages_sim", 0) == 0:
            penalty = 0.8
            
        final_score = score * penalty

        swipes.append({
            "client_id": client.get("profile_id"),
            "therapist_id": therapist_id,
            "points": round(final_score, 2),
            "reasons": reasons,
            "status": "pending",
            "created_at": now
        })

    swipes.sort(key=lambda x: x["points"], reverse=True)
    return swipes
