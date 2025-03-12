# URL Shortening Service

A simple URL shortener built with Go, Gin, and MySQL.

## System Design Features (Monorepo)
- [x] Shorten URLs with custom aliases
- [x] Redirect shortened URLs
- [x] Rate limiting to prevent abuse
- [x] Track URL clicks (updates)
- [x] Implement structured logging
- [x] Implement a ID generation mechanism 
- [x] Timestamp of creation and last access
- [ ] Implement Redis/Memcached for caching frequent URL redirects
- [ ] Create a separate analytics microservice
- [ ] Soft delete mechanism for unused/expired URLs
- [ ] Periodic cleanup job to remove expired entries
- [ ] Block URLs from known malicious domains
- [x] Implement user registration and authentication
- [ ] OAuth/JWT-based authentication
- [x] Role-based access control (RBAC)
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

## Test URL
- **Method**: GET  
- **Endpoint**: `/ping`

## Redirect
- **Method**: GET  
- **Endpoint**: `/0FL0EApfD`

## Short URL Create
- **Method**: POST  
- **Endpoint**: `/api/v1`  
- **Headers**: `Content-Type: application/json`  
- **Body**:
  ```json
  {
     "url": "http://example.com",
     "short": "abc",
     "expiry": 24
  }
  ```

## Get Analytics
- **Method**: GET  
- **Endpoint**: `/api/v1/analytics/0FL0EApfD`  
- **Headers**: `Content-Type: application/json`

## Users Register
- **Method**: POST  
- **Endpoint**: `/api/v1/register`  
- **Headers**: `Content-Type: application/json`  
- **Body**:
  ```json
  {
     "username": "amir",
     "password": "Qwerty"
  }
  ```

## Users Login
- **Method**: POST  
- **Endpoint**: `/api/v1/login`  
- **Headers**: `Content-Type: application/json`  
- **Body**:
  ```json
  {
     "username": "amir",
     "password": "Qwerty"
  }
  ```


## License
MIT License.
