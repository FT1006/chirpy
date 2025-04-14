# Chirpy

Chirpy is a backend API service for a Twitter-like social platform where users can post short messages called "chirps". Built with Go, it provides user authentication, data persistence, and premium membership features.

## How to Use

### Prerequisites
- Go 1.16 or newer
- PostgreSQL
- Environment variables properly configured

### Setup and Installation
1. Clone the repository:
   ```
   git clone https://github.com/yourusername/chirpy.git
   cd chirpy
   ```

2. Set up environment variables:
   ```
   # Create a .env file
   JWT_SECRET=your_secret_key
   DATABASE_URL=postgres://username:password@localhost:5432/chirpy
   POLKA_API_KEY=your_api_key
   ```

3. Initialize the database:
   ```
   # Run PostgreSQL migrations
   goose -dir sql/schema postgres "postgres://username:password@localhost:5432/chirpy" up
   ```

4. Build and run:
   ```
   go build
   ./chirpy
   ```

### Quick Start Guide
1. **Create a user account**:
   ```
   curl -X POST http://localhost:8080/api/users \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"securepassword"}'
   ```

2. **Login**:
   ```
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"securepassword"}'
   ```
   Save the returned JWT token and refresh token.

3. **Create a chirp**:
   ```
   curl -X POST http://localhost:8080/api/chirps \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"body":"Hello, Chirpy!"}'
   ```

4. **Get all chirps**:
   ```
   curl http://localhost:8080/api/chirps
   ```

## Technical Architecture

Chirpy is built using:
- Go's standard library `net/http` package
- PostgreSQL database for data persistence
- JWT authentication system
- RESTful API design

## Database Structure

The PostgreSQL database consists of the following tables:

1. Users
   - UUID primary key
   - Email (unique)
   - Hashed password
   - Created/updated timestamps
   - Chirpy Red membership status

2. Chirps
   - UUID primary key
   - Message content (max 140 characters)
   - User ID (foreign key)
   - Created/updated timestamps

3. Refresh Tokens
   - Token string primary key
   - User ID foreign key
   - Created/updated/expiry timestamps
   - Revoked timestamp for invalidation

## Authentication System

- Passwords hashed using bcrypt
- JWT tokens with 1-hour expiry
- Refresh token system for extended sessions
- Bearer token authentication for protected endpoints
- API key validation for webhook endpoints

## API Endpoints

### User Management

#### Create User
- **Endpoint**: `POST /api/users`
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "email": "string",
    "is_chirpy_red": boolean
  }
  ```

#### Login
- **Endpoint**: `POST /api/login`
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "email": "string",
    "token": "string",
    "refresh_token": "string",
    "is_chirpy_red": boolean
  }
  ```

#### Update User
- **Endpoint**: `PUT /api/users`
- **Auth**: Bearer token
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "email": "string",
    "is_chirpy_red": boolean
  }
  ```

#### Refresh Token
- **Endpoint**: `POST /api/refresh`
- **Auth**: Bearer token (refresh token)
- **Response** (200 OK):
  ```json
  {
    "token": "string"
  }
  ```

#### Revoke Token
- **Endpoint**: `POST /api/revoke`
- **Auth**: Bearer token (refresh token)
- **Response**: 204 No Content

### Chirp Operations

#### Create Chirp
- **Endpoint**: `POST /api/chirps`
- **Auth**: Bearer token
- **Request Body**:
  ```json
  {
    "body": "string" // Max 140 characters
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "body": "string",
    "user_id": "uuid"
  }
  ```
- **Note**: Profanity filter replaces words like "kerfuffle", "sharbert", "fornax" with "****"

#### Get All Chirps
- **Endpoint**: `GET /api/chirps`
- **Query Parameters**:
  - `author_id`: (optional) Filter by author UUID
  - `sort`: (optional) "asc" or "desc" to sort by creation time
- **Response** (200 OK):
  ```json
  [
    {
      "id": "uuid",
      "created_at": "timestamp",
      "updated_at": "timestamp",
      "body": "string",
      "user_id": "uuid"
    }
  ]
  ```

#### Get Chirp by ID
- **Endpoint**: `GET /api/chirps/{chirpID}`
- **Response** (200 OK):
  ```json
  {
    "id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "body": "string",
    "user_id": "uuid"
  }
  ```

#### Delete Chirp
- **Endpoint**: `DELETE /api/chirps/{chirpID}`
- **Auth**: Bearer token (must be creator of chirp)
- **Response**: 204 No Content

### System Endpoints

#### Health Check
- **Endpoint**: `GET /api/healthz`
- **Response**: 200 OK

#### Metrics
- **Endpoint**: `GET /admin/metrics`
- **Response**: Site metrics

#### Reset
- **Endpoint**: `POST /admin/reset`
- **Response**: 200 OK

#### Polka Webhook
- **Endpoint**: `POST /api/polka/webhooks`
- **Auth**: API key
- **Request Body**:
  ```json
  {
    "event": "string",
    "data": {
      "user_id": "string"
    }
  }
  ```
- **Response**: 204 No Content
- **Note**: Handles "user.upgraded" event to update premium membership status

## Features

- Profanity Filter: Automatically replaces forbidden words with asterisks
- Chirpy Red: Premium membership tier activated through Polka webhooks
- Metrics: Tracking of application access
- Secure Token Management: JWT creation and validation
- Data Persistence: All data stored in PostgreSQL

## Dependencies

- PostgreSQL
- JWT
- bcrypt
- UUID
- godotenv

## Development Tools

- sqlc for SQL query generation
- goose for database migrations