# Blog Service - PsyConnect

## Overview
The **Blog Service** manages blog posts, articles, and educational content within the PsyConnect platform. This service provides mental health resources, educational materials, and community-generated content to help users on their mental health journey.

## Features
- Blog post creation and management
- Article categorization and tagging
- Content moderation and approval
- Comment system and user engagement
- Search and filtering capabilities
- SEO optimization for content

## Technology Stack
- **Backend**: [Technology to be implemented]
- **Database**: [Database to be implemented]
- **Search Engine**: Elasticsearch (planned)
- **API Communication**: RESTful APIs
- **Content Management**: Rich text editor support

## API Endpoints

### Blog Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/blog` | Get all blog posts (paginated) |
| POST | `/blog` | Create new blog post |
| GET | `/blog/{id}` | Get specific blog post |
| PUT | `/blog/{id}` | Update blog post |
| DELETE | `/blog/{id}` | Delete blog post |

### Category Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/blog/categories` | Get all categories |
| POST | `/blog/categories` | Create new category |
| GET | `/blog/category/{id}` | Get posts by category |
| PUT | `/blog/categories/{id}` | Update category |
| DELETE | `/blog/categories/{id}` | Delete category |

### Comment Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/blog/{id}/comments` | Get comments for blog post |
| POST | `/blog/{id}/comments` | Add comment to blog post |
| PUT | `/blog/{id}/comments/{commentId}` | Update comment |
| DELETE | `/blog/{id}/comments/{commentId}` | Delete comment |

### Search and Filtering
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/blog/search` | Search blog posts |
| GET | `/blog/tags/{tag}` | Get posts by tag |
| GET | `/blog/author/{authorId}` | Get posts by author |

*Note: This service is currently in development. API endpoints are subject to change.*

## Setup & Configuration

### Environment Variables
```env
PORT=8085
DATABASE_URL={your-database-url}
ELASTICSEARCH_URL={your-elasticsearch-url}
JWT_SECRET={your-jwt-secret}
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/blogservice
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
