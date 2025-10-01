# Inter-Service Communication Test Guide

This document demonstrates how to test the inter-service communication between `open_data_service` and `st_dom_service`.

## Prerequisites

1. MongoDB running on `localhost:27018`
2. Both services compiled and ready to run

## Starting the Services

### 1. Start st_dom_service (Port 8081)
```bash
cd st_dom_service
./st_dom_service
```

### 2. Start open_data_service (Port 8082)
```bash
cd open_data_service
./open_data_service
```

## Testing Inter-Service Communication

### 1. Check st_dom_service health
```bash
curl http://localhost:8081/health
```

### 2. Check open_data_service health
```bash
curl http://localhost:8082/health
```

### 3. Test inter-service health check
```bash
curl http://localhost:8082/api/v1/inter-service/health
```

### 4. Get all accepted applications via inter-service communication
```bash
# Without authentication (if endpoint allows)
curl http://localhost:8082/api/v1/inter-service/prihvacene-aplikacije

# With authentication (Authorization header is forwarded to st_dom_service)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8082/api/v1/inter-service/prihvacene-aplikacije
```

### 5. Get accepted applications for a specific user
```bash
# With authentication (Authorization header is forwarded to st_dom_service)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8082/api/v1/inter-service/prihvacene-aplikacije/user/{USER_ID}
```

### 6. Get accepted applications for a specific room
```bash
# With authentication (Authorization header is forwarded to st_dom_service)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8082/api/v1/inter-service/prihvacene-aplikacije/room/{ROOM_ID}
```

### 7. Get accepted applications for a specific academic year
```bash
# With authentication (Authorization header is forwarded to st_dom_service)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8082/api/v1/inter-service/prihvacene-aplikacije/academic-year/2024/2025
```

## Expected Response Format

All inter-service endpoints will return responses in the following format:

```json
{
  "message": "Data retrieved from st_dom_service via inter-service communication",
  "data": [
    {
      "id": "application_id",
      "aplikacija_id": "original_application_id",
      "user_id": "user_id",
      "broj_indexa": "student_index",
      "prosek": 9,
      "soba_id": "room_id",
      "academic_year": "2024/2025",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "count": 1,
  "source": "st_dom_service"
}
```

## Error Handling

If st_dom_service is unavailable, the inter-service endpoints will return:

```json
{
  "error": "Failed to retrieve accepted applications from st_dom_service: [error details]"
}
```

## Architecture Overview

```
Client Request → open_data_service → st_dom_service → Response
     ↓                ↓                    ↓            ↓
   Port 8082      HTTP Client        Port 8081    JSON Response
   (with JWT)     (forwards JWT)     (validates JWT)
```

The `open_data_service` acts as a proxy/gateway that:
1. Receives requests from clients (with optional Authorization header)
2. Extracts the Authorization header from the incoming request
3. Makes HTTP calls to `st_dom_service` with the forwarded Authorization header
4. Forwards the response back to the client
5. Adds metadata about the inter-service communication

## Authorization Header Forwarding

The inter-service communication now supports forwarding the Authorization header:

- **Client** sends request to `open_data_service` with `Authorization: Bearer <token>`
- **open_data_service** extracts the header and forwards it to `st_dom_service`
- **st_dom_service** validates the JWT token and processes the request
- **Response** is returned through the same chain

This ensures that:
- User authentication is preserved across service boundaries
- Authorization policies in `st_dom_service` are respected
- The same JWT token works for both direct and inter-service calls

This demonstrates a microservices architecture where services communicate over HTTP APIs with proper authentication forwarding.
