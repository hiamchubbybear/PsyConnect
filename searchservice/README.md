# Search Service - PsyConnect

## Overview
The **Search Service** provides comprehensive search functionality across the PsyConnect platform. This service enables users to search for therapists, content, resources, and other platform data efficiently and accurately.

## Features
- Global search across all platform content
- Advanced filtering and sorting options
- Therapist search with specialization filters
- Content search (blogs, articles, resources)
- Auto-complete and suggestions
- Search analytics and optimization

## Technology Stack
- **Backend**: [Technology to be implemented]
- **Search Engine**: Elasticsearch
- **Database**: [Database to be implemented]
- **API Communication**: RESTful APIs
- **Caching**: Redis

## API Endpoints

### General Search
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/search` | Global search across platform |
| GET | `/search/suggestions` | Get search suggestions |
| GET | `/search/autocomplete` | Auto-complete search queries |

### Therapist Search
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/search/therapists` | Search therapists |
| GET | `/search/therapists/filters` | Get available search filters |
| GET | `/search/therapists/specializations` | Get therapist specializations |
| GET | `/search/therapists/location` | Search therapists by location |

### Content Search
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/search/content` | Search blog posts and articles |
| GET | `/search/resources` | Search mental health resources |
| GET | `/search/categories` | Search content by categories |

### Advanced Search
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/search/advanced` | Advanced search with complex filters |
| GET | `/search/facets` | Get faceted search results |
| GET | `/search/trending` | Get trending search terms |

### Search Analytics
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/search/analytics` | Get search analytics data |
| POST | `/search/track` | Track search queries |

*Note: This service is currently in development. API endpoints are subject to change.*

## Search Parameters

### Common Query Parameters
- `q` - Search query string
- `limit` - Number of results per page (default: 20)
- `offset` - Pagination offset (default: 0)
- `sort` - Sort order (relevance, date, popularity)
- `filter` - Filter criteria

### Therapist Search Filters
- `specialization` - Therapist specialization
- `location` - Geographic location
- `availability` - Available time slots
- `rating` - Minimum rating
- `price_range` - Price range filter

## Setup & Configuration

### Environment Variables
```env
PORT=8086
ELASTICSEARCH_URL={your-elasticsearch-url}
REDIS_URL={your-redis-url}
DATABASE_URL={your-database-url}
JWT_SECRET={your-jwt-secret}
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/searchservice
   ```
2. Install dependencies and start the service (commands will be added when implemented)

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
