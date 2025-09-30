# Postman Testing Guide - Open Data Service

This guide provides all endpoints with example requests and responses for testing in Postman.

**Base URL:** `http://localhost:8082`

---

## 1. Health Check

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/health`
- **Query Parameters:** None
- **Body:** None

### Example Response (200 OK)
```json
{
  "status": "ok",
  "service": "open_data_service"
}
```

---

## 2. Get All Rooms

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms`
- **Query Parameters:** None
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439013",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 3,
      "luksuzi": ["klima", "sopstveno kupatilo"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 2
}
```

---

## 3. Get All Student Dormitories

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/st-doms`
- **Query Parameters:** None
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439012",
      "address": "Studentska 123, Beograd",
      "telephone_number": "+381 11 1234567",
      "email": "dorm1@example.com",
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439014",
      "address": "Bulevar Revolucije 45, Beograd",
      "telephone_number": "+381 11 7654321",
      "email": "dorm2@example.com",
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 2
}
```

---

## 4. Filter Rooms by Luxury Amenities (Single)

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/filter-by-luksuz`
- **Query Parameters:**
  - `luksuzi` = `klima`
- **Full URL:** `http://localhost:8082/api/v1/rooms/filter-by-luksuz?luksuzi=klima`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439013",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 3,
      "luksuzi": ["klima", "sopstveno kupatilo"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 2
}
```

---

## 5. Filter Rooms by Luxury Amenities (Multiple)

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/filter-by-luksuz`
- **Query Parameters:**
  - `luksuzi` = `klima,terasa`
- **Full URL:** `http://localhost:8082/api/v1/rooms/filter-by-luksuz?luksuzi=klima,terasa`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 1
}
```

### Error Response (400 Bad Request) - Invalid Luxury
```json
{
  "error": "Invalid luxury amenity: invalid_amenity"
}
```

---

## 6. Filter Rooms by Luxury Amenities (All Valid Values)

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/filter-by-luksuz`
- **Query Parameters:**
  - `luksuzi` = `klima,terasa,sopstveno kupatilo`
- **Full URL:** `http://localhost:8082/api/v1/rooms/filter-by-luksuz?luksuzi=klima,terasa,sopstveno%20kupatilo`
- **Body:** None
- **Note:** URL encode spaces as `%20`

### Valid Luxury Amenities
- `klima` - Air conditioning
- `terasa` - Terrace
- `sopstveno kupatilo` - Private bathroom (use `sopstveno%20kupatilo` in URL)
- `áram` - Electric current
- `ablak` - Window
- `neisvrljan zid` - Unpainted wall (use `neisvrljan%20zid` in URL)

---

## 7. Filter Rooms by Luxury and Student Dormitory

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/filter-by-luksuz-and-stdom`
- **Query Parameters:**
  - `st_dom_id` = `507f1f77bcf86cd799439012` (required)
  - `luksuzi` = `klima` (optional)
- **Full URL:** `http://localhost:8082/api/v1/rooms/filter-by-luksuz-and-stdom?st_dom_id=507f1f77bcf86cd799439012&luksuzi=klima`
- **Body:** None
- **Note:** Replace `st_dom_id` with actual ObjectID from your database

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z",
      "st_dom": {
        "id": "507f1f77bcf86cd799439012",
        "address": "Studentska 123, Beograd",
        "telephone_number": "+381 11 1234567",
        "email": "dorm1@example.com",
        "created_at": "2025-09-30T10:00:00Z",
        "updated_at": "2025-09-30T10:00:00Z"
      }
    }
  ],
  "count": 1
}
```

### Error Response (400 Bad Request) - Missing st_dom_id
```json
{
  "error": "st_dom_id is required"
}
```

### Error Response (400 Bad Request) - Invalid st_dom_id format
```json
{
  "error": "Invalid st_dom_id format"
}
```

---

## 8. Filter Rooms by Exact Bed Capacity

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/filter-by-krevetnost`
- **Query Parameters:**
  - `exact` = `2`
- **Full URL:** `http://localhost:8082/api/v1/rooms/filter-by-krevetnost?exact=2`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 1
}
```

---

## 9. Filter Rooms by Bed Capacity Range

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/filter-by-krevetnost`
- **Query Parameters:**
  - `min` = `2`
  - `max` = `4`
- **Full URL:** `http://localhost:8082/api/v1/rooms/filter-by-krevetnost?min=2&max=4`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439013",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 3,
      "luksuzi": ["klima", "sopstveno kupatilo"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 2
}
```

### Error Response (400 Bad Request) - Invalid range
```json
{
  "error": "min value cannot be greater than max value"
}
```

---

## 10. Filter Rooms by Minimum Bed Capacity

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/filter-by-krevetnost`
- **Query Parameters:**
  - `min` = `3`
- **Full URL:** `http://localhost:8082/api/v1/rooms/filter-by-krevetnost?min=3`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439013",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 3,
      "luksuzi": ["klima", "sopstveno kupatilo"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439015",
      "st_dom_id": "507f1f77bcf86cd799439014",
      "krevetnost": 4,
      "luksuzi": ["klima"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 2
}
```

---

## 11. Filter Rooms by Maximum Bed Capacity

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/filter-by-krevetnost`
- **Query Parameters:**
  - `max` = `2`
- **Full URL:** `http://localhost:8082/api/v1/rooms/filter-by-krevetnost?max=2`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439016",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 1,
      "luksuzi": ["klima", "sopstveno kupatilo"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 2
}
```

---

## 12. Search Dormitories by Address (Simple)

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/st-doms/search-by-address`
- **Query Parameters:**
  - `address` = `Studentska`
- **Full URL:** `http://localhost:8082/api/v1/st-doms/search-by-address?address=Studentska`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439012",
      "address": "Studentska 123, Beograd",
      "telephone_number": "+381 11 1234567",
      "email": "dorm1@example.com",
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 1
}
```

### Error Response (400 Bad Request) - Missing address
```json
{
  "error": "address parameter is required"
}
```

---

## 13. Search Dormitories by Address (Regex - Starts With)

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/st-doms/search-by-address`
- **Query Parameters:**
  - `address` = `^Bul`
- **Full URL:** `http://localhost:8082/api/v1/st-doms/search-by-address?address=^Bul`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439014",
      "address": "Bulevar Revolucije 45, Beograd",
      "telephone_number": "+381 11 7654321",
      "email": "dorm2@example.com",
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 1
}
```

---

## 14. Search Dormitories by Address (Contains "Beograd")

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/st-doms/search-by-address`
- **Query Parameters:**
  - `address` = `Beograd`
- **Full URL:** `http://localhost:8082/api/v1/st-doms/search-by-address?address=Beograd`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439012",
      "address": "Studentska 123, Beograd",
      "telephone_number": "+381 11 1234567",
      "email": "dorm1@example.com",
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439014",
      "address": "Bulevar Revolucije 45, Beograd",
      "telephone_number": "+381 11 7654321",
      "email": "dorm2@example.com",
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z"
    }
  ],
  "count": 2
}
```

---

## 15. Advanced Filter (All Parameters)

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/advanced-filter`
- **Query Parameters:**
  - `luksuzi` = `klima,terasa`
  - `st_dom_id` = `507f1f77bcf86cd799439012`
  - `min` = `2`
  - `max` = `4`
- **Full URL:** `http://localhost:8082/api/v1/rooms/advanced-filter?luksuzi=klima,terasa&st_dom_id=507f1f77bcf86cd799439012&min=2&max=4`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z",
      "st_dom": {
        "id": "507f1f77bcf86cd799439012",
        "address": "Studentska 123, Beograd",
        "telephone_number": "+381 11 1234567",
        "email": "dorm1@example.com",
        "created_at": "2025-09-30T10:00:00Z",
        "updated_at": "2025-09-30T10:00:00Z"
      }
    }
  ],
  "count": 1
}
```

---

## 16. Advanced Filter (Luxury and Exact Bed Capacity)

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/advanced-filter`
- **Query Parameters:**
  - `luksuzi` = `klima`
  - `exact` = `2`
- **Full URL:** `http://localhost:8082/api/v1/rooms/advanced-filter?luksuzi=klima&exact=2`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z",
      "st_dom": {
        "id": "507f1f77bcf86cd799439012",
        "address": "Studentska 123, Beograd",
        "telephone_number": "+381 11 1234567",
        "email": "dorm1@example.com",
        "created_at": "2025-09-30T10:00:00Z",
        "updated_at": "2025-09-30T10:00:00Z"
      }
    }
  ],
  "count": 1
}
```

---

## 17. Advanced Filter (Only Dormitory and Bed Range)

### Request
- **Method:** `GET`
- **URL:** `http://localhost:8082/api/v1/rooms/advanced-filter`
- **Query Parameters:**
  - `st_dom_id` = `507f1f77bcf86cd799439012`
  - `min` = `2`
- **Full URL:** `http://localhost:8082/api/v1/rooms/advanced-filter?st_dom_id=507f1f77bcf86cd799439012&min=2`
- **Body:** None

### Example Response (200 OK)
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 2,
      "luksuzi": ["klima", "terasa"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z",
      "st_dom": {
        "id": "507f1f77bcf86cd799439012",
        "address": "Studentska 123, Beograd",
        "telephone_number": "+381 11 1234567",
        "email": "dorm1@example.com",
        "created_at": "2025-09-30T10:00:00Z",
        "updated_at": "2025-09-30T10:00:00Z"
      }
    },
    {
      "id": "507f1f77bcf86cd799439013",
      "st_dom_id": "507f1f77bcf86cd799439012",
      "krevetnost": 3,
      "luksuzi": ["klima", "sopstveno kupatilo"],
      "created_at": "2025-09-30T10:00:00Z",
      "updated_at": "2025-09-30T10:00:00Z",
      "st_dom": {
        "id": "507f1f77bcf86cd799439012",
        "address": "Studentska 123, Beograd",
        "telephone_number": "+381 11 1234567",
        "email": "dorm1@example.com",
        "created_at": "2025-09-30T10:00:00Z",
        "updated_at": "2025-09-30T10:00:00Z"
      }
    }
  ],
  "count": 2
}
```

---

## Postman Collection Setup

### Quick Import
1. Create a new Collection in Postman called "Open Data Service"
2. Set the collection variable `baseUrl` = `http://localhost:8082`
3. Add each request above as a separate request in the collection

### Environment Variables
Create a Postman environment with these variables:
- `baseUrl`: `http://localhost:8082`
- `st_dom_id`: `507f1f77bcf86cd799439012` (replace with actual ID from your database)

### Tips for Testing
1. **Start with health check** to ensure the service is running
2. **Get all rooms/dormitories** first to obtain real IDs from your database
3. **Copy an actual ObjectID** from the response to use in `st_dom_id` parameter
4. **Test error cases** by providing invalid values
5. **Try different combinations** of parameters

---

## Common Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid luxury amenity: xyz"
}
```

### 400 Bad Request
```json
{
  "error": "Invalid st_dom_id format"
}
```

### 400 Bad Request
```json
{
  "error": "min value cannot be greater than max value"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to filter rooms"
}
```

---

## Summary

All endpoints use **GET** method with **query parameters** (no JSON body required).

**Total Endpoints:** 8
1. `/health` - Health check
2. `/api/v1/rooms` - Get all rooms
3. `/api/v1/st-doms` - Get all dormitories
4. `/api/v1/rooms/filter-by-luksuz` - Filter by luxury
5. `/api/v1/rooms/filter-by-luksuz-and-stdom` - Filter by luxury + dorm
6. `/api/v1/rooms/filter-by-krevetnost` - Filter by bed capacity
7. `/api/v1/st-doms/search-by-address` - Search dormitories by address
8. `/api/v1/rooms/advanced-filter` - Combined multi-criteria filter

**No authentication required** - all endpoints are public!

