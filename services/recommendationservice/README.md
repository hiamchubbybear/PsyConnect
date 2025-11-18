# Recommendation Service - PsyConnect

## Overview
The **Recommendation Service** provides AI-powered recommendations for matching clients with therapists within the PsyConnect platform. This service uses machine learning algorithms to provide personalized therapist recommendations based on client preferences, needs, and compatibility factors.

## Features
- AI-powered therapist-client matching
- Similarity-based scoring algorithms
- Multi-criteria recommendation system
- Real-time matching capabilities
- Data-driven insights and analytics
- Continuous learning and improvement

## Technology Stack
- **Backend**: Python
- **Machine Learning**: scikit-learn, pandas, numpy
- **Algorithms**: Cosine similarity, collaborative filtering
- **Data Format**: JSON
- **API Communication**: RESTful APIs (planned)

## AI Algorithms

### Cosine Similarity Matching
The service uses cosine similarity to match clients with therapists based on:
- Specialization alignment
- Language compatibility  
- Location preferences
- Availability matching
- Price range compatibility
- Experience level requirements

### Matching Criteria
- **Specialization**: Client issues vs therapist expertise
- **Languages**: Communication language preferences
- **Location**: Geographic proximity or online availability
- **Price Range**: Budget compatibility
- **Availability**: Schedule alignment
- **Experience Level**: Required vs available experience
- **Gender Preferences**: Optional therapist gender preference
- **Session Duration**: Preferred session length

## API Endpoints

### Recommendation Engine
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/recommend/therapist` | Get therapist recommendations for client |
| GET | `/recommend/match/{clientId}` | Get pre-computed matches for client |
| POST | `/recommend/update` | Update recommendation parameters |

### Analytics and Insights
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/analytics/match-success` | Get matching success statistics |
| GET | `/analytics/popular-specializations` | Get popular specialization trends |
| POST | `/analytics/feedback` | Submit matching feedback |

### Data Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/data/therapist` | Add/update therapist data |
| POST | `/data/client` | Add/update client preferences |
| GET | `/data/statistics` | Get system statistics |

*Note: API endpoints are currently in development.*

## Matching Process

### Input Parameters
```json
{
  "profile_id": "C002",
  "address": "HCMC",
  "languages": ["English"],
  "issue_detail": "marriage",
  "consultation_modes": ["online"],
  "range_price": 600000,
  "availability": ["Morning_Weekday"],
  "preferred_therapist_gender": null,
  "experience_level": "senior",
  "specialization": ["marriage"],
  "urgency_level": "low",
  "session_duration": 60,
  "preferred_therapist_language": "English",
  "is_flexible_with_schedule": false
}
```

### Output Format
```json
{
  "matches": [
    {
      "therapist_id": "T002",
      "compatibility_score": 0.95,
      "match_reasons": ["specialization_match", "language_match", "location_match"],
      "confidence": "high"
    }
  ],
  "total_matches": 5,
  "processing_time_ms": 45
}
```

## Data Files
- `generated_therapists.json` - Therapist profile data
- `generated_client.json` - Client preference data  
- `matching_results.json` - Pre-computed matching results
- `hard_client.json` - Test client data

## Setup & Configuration

### Environment Variables
```env
PORT=8087
MODEL_PATH={path-to-model-files}
DATA_PATH={path-to-data-files}
LOG_LEVEL=INFO
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/recommendationservice
   ```
2. Install Python dependencies:
   ```bash
   pip install pandas numpy scikit-learn
   ```
3. Run the recommendation engine:
   ```bash
   python recommendation-ai/consine.py
   ```

## Performance
- Processing time: ~45ms per recommendation request
- Supports real-time matching for up to 1000+ therapists
- Scalable architecture for growing user base

## Contributing
We welcome contributions! Please follow the standard Git workflow:
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/YourFeature`)
3. Commit your changes (`git commit -m 'Add YourFeature'`)
4. Push to your branch (`git push origin feature/YourFeature`)
5. Open a Pull Request

## Contact
For inquiries, reach out via:
- **Project Lead**: Chessy
- **Email**: [tranvanhuy16032004@gmail.com](mailto:tranvanhuy16032004@gmail.com)
- **GitHub Repository**: [PsyConnect](https://github.com/hiamchubbybear/PsyConnect)
