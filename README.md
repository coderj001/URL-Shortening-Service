# URL Shortening Service

A simple URL shortener built with Go, Gin, and MySQL.

## Features
- Shorten URLs with custom aliases
- Redirect shortened URLs
- Rate limiting to prevent abuse
- Docker support for easy deployment

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
