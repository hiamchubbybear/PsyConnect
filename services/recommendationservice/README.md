# Recommendation Service - PsyConnect

![Python](https://img.shields.io/badge/Python-3.9-blue?style=flat&logo=python&logoColor=white)
![Flask](https://img.shields.io/badge/Flask-2.0-black?style=flat&logo=flask&logoColor=white)
![Scikit-learn](https://img.shields.io/badge/Scikit_Learn-ML-orange?style=flat&logo=scikitlearn&logoColor=white)

## 📖 Overview

The **Recommendation Service** powers the smart matching between clients and therapists. It uses Machine Learning algorithms (Cosine Similarity) to analyze profiles and suggest the best matches based on specialization, language, availability, and more.

## ✨ Features

- **Smart Matching**: Calculates compatibility scores between clients and therapists.
- **Multi-criteria Analysis**: Considers language, location, price, and expertise.
- **Real-time Recommendations**: Provides instant suggestions via API.

## 🛠 Technology Stack

- **Language**: Python 3.9
- **Framework**: Flask
- **ML Libraries**: Scikit-learn, Pandas, NumPy
- **Algorithm**: Cosine Similarity, Collaborative Filtering

## 🔌 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/recommend` | Get therapist recommendations for a client |
| `GET` | `/health` | Service health check |

## ⚙️ Configuration

### Environment Variables
This service uses a `.env` file for configuration.
1. Copy the example file:
   ```bash
   cp .env.example .env
   ```
2. Update the variables in `.env`:
   - `PORT`: Service port (default: 8086)
   - `FLASK_ENV`: development/production

## 🚀 Installation & Run

### Prerequisites
- Python 3.9+
- Pip

### Local Run
```bash
# 1. Navigate to directory
cd services/recommendationservice

# 2. Install dependencies
pip install -r requirements.txt

# 3. Run application
python api/main.py
```

### Docker Run
```bash
docker build -t psyconnect/recommendationservice .
docker run -p 8086:8086 --env-file .env psyconnect/recommendationservice
```

## 🤝 Contributing
Please refer to the root [README](../../README.md) for contributing guidelines.
