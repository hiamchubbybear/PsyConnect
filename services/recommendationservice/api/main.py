from flask import Flask, request, jsonify
from sklearn.preprocessing import MultiLabelBinarizer
from sklearn.metrics.pairwise import cosine_similarity
import pandas as pd
from datetime import datetime

app = Flask(__name__)

@app.route('/recommend', methods=['POST'])
def recommend():
    try:
        data = request.json
        client = data.get("clientRaw")
        therapists = data.get("therapistsRaw")

        if not client or not therapists:
            return jsonify({"error": "Missing client or therapist data"}), 400

        from services.matching import calculate_similarity
        swipes = calculate_similarity(client, therapists)

        result = {
            "client_id": client.get("profile_id"),  
            "swipes": swipes
        }
        return jsonify(result)

    except Exception as e:
        print(" Internal error:", e)
        return jsonify({"error": "Internal server error", "details": str(e)}), 500

@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "UP"}), 200

import os

if __name__ == '__main__':
    port = int(os.environ.get("PORT", 8086))
    app.run(host="0.0.0.0", port=port, debug=True)
