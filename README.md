# URL Shortening Service

A simple URL shortener built with Go, Gin, and MySQL.

## System Design Features (Monorepo)
- [x] Shorten URLs with custom aliases
- [x] Redirect shortened URLs
- [x] Rate limiting to prevent abuse
- [x] Track URL clicks
- [ ] Implement structured logging
- [ ] Implement a ID generation mechanism 
- [ ] Timestamp of creation and last access
- [ ] Implement Redis/Memcached for caching frequent URL redirects
- [ ] Create a separate analytics microservice
- [ ] Soft delete mechanism for unused/expired URLs
- [ ] Periodic cleanup job to remove expired entries
- [ ] Block URLs from known malicious domains
- [ ] Implement user registration and authentication
- [ ] OAuth/JWT-based authentication
- [ ] Role-based access control (RBAC)
  - [ ] Free tier users
  - [ ] Premium users with additional features
- [ ] Optimize database queries and indexing
- [ ] Add TDD 
- [ ] Integration Tests
- [ ] Frontend/Template HTmX

## Prerequisites
- Docker
- Docker Compose

## Running with Docker

1. Build and start the services:
   ```sh
   docker-compose up --build
   ```
2. The service will be available at `http://localhost:3000`.

## API Endpoints
### Shorten a URL
- **Endpoint:** `POST /api/v1`
- **Request Body:**
  ```json
  {
    "url": "https://example.com",
    "short": "customAlias",  // Optional
    "expiry": 24  // Expiry in hours (Optional, default: 24h)
  }
  ```
- **Response:**
  ```json
  {
    "url": "https://example.com",
    "short": "localhost:3000/customAlias",
    "expiry": 24,
    "rate_limit": 10,
    "rate_limit_reset": 30
  }
  ```

### Resolve a Short URL
- **Endpoint:** `GET /:shortURL`
- **Response:** Redirects to the original URL.



## License
MIT License.