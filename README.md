# SSO Service

A Go microservice for user authentication and authorization using MongoDB and JWT tokens.

## Features

- User registration with validation
- User login with JWT token generation
- Password hashing using bcrypt
- MongoDB integration
- CORS support
- Role-based access control
- Protected routes with JWT middleware

## User Model

The service manages users with the following fields:
- `username` - Unique username (3-20 characters)
- `email` - Valid email address
- `password` - Securely hashed password (minimum 6 characters)
- `first_name` - User's first name
- `last_name` - User's last name
- `role` - User role (default: "user")

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### Protected Endpoints (Require JWT Token)

- `GET /api/v1/profile` - Get current user profile

## Configuration

Copy `config.env` and update the values:

```env
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=sso_db
JWT_SECRET=your_jwt_secret_key_here_change_in_production
PORT=8080
GIN_MODE=debug
```

## Installation and Setup

1. Install Go dependencies:
```bash
go mod tidy
```

2. Make sure MongoDB is running on your system

3. Update the configuration in `config.env`

4. Run the service:
```bash
go run main.go
```

## API Usage Examples

### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get user profile (requires JWT token)
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## Project Structure

```
sso_service/
├── config/          # Configuration management
├── database/        # MongoDB connection
├── handlers/        # HTTP request handlers
├── middleware/      # Authentication and CORS middleware
├── models/          # Data models and structs
├── routes/          # Route definitions
├── services/        # Business logic
├── utils/           # Utility functions (JWT, password hashing)
├── main.go          # Application entry point
├── go.mod           # Go module file
└── config.env       # Environment configuration
```

## Security Features

- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- CORS middleware for cross-origin requests
- Input validation on all endpoints
- Protected routes require valid JWT tokens

## Dependencies

- `gin-gonic/gin` - HTTP web framework
- `go.mongodb.org/mongo-driver` - MongoDB driver
- `golang-jwt/jwt/v4` - JWT token handling
- `golang.org/x/crypto` - Password hashing
- `joho/godotenv` - Environment variable loading

