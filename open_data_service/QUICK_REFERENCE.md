# Open Data Service - Quick Reference

**Base URL:** `http://localhost:8082`

## All Endpoints at a Glance

| # | Endpoint | Method | Query Parameters | Description |
|---|----------|--------|------------------|-------------|
| 1 | `/health` | GET | None | Health check |
| 2 | `/api/v1/rooms` | GET | None | Get all rooms |
| 3 | `/api/v1/st-doms` | GET | None | Get all dormitories |
| 4 | `/api/v1/rooms/filter-by-luksuz` | GET | `luksuzi` (optional) | Filter by luxury amenities |
| 5 | `/api/v1/rooms/filter-by-luksuz-and-stdom` | GET | `luksuzi` (optional), `st_dom_id` (required) | Filter by luxury + dorm |
| 6 | `/api/v1/rooms/filter-by-krevetnost` | GET | `exact`, `min`, `max` (all optional) | Filter by bed capacity |
| 7 | `/api/v1/st-doms/search-by-address` | GET | `address` (required) | Search dorms by address |
| 8 | `/api/v1/rooms/advanced-filter` | GET | `luksuzi`, `st_dom_id`, `exact`, `min`, `max` (all optional) | Combined multi-filter |

---

## Copy-Paste URLs for Postman

### 1. Health Check
```
GET http://localhost:8082/health
```

### 2. Get All Rooms
```
GET http://localhost:8082/api/v1/rooms
```

### 3. Get All Student Dormitories
```
GET http://localhost:8082/api/v1/st-doms
```

### 4. Filter by Single Luxury Amenity
```
GET http://localhost:8082/api/v1/rooms/filter-by-luksuz?luksuzi=klima
```

### 5. Filter by Multiple Luxury Amenities
```
GET http://localhost:8082/api/v1/rooms/filter-by-luksuz?luksuzi=klima,terasa
```

### 6. Filter by Luxury + Dormitory
```
GET http://localhost:8082/api/v1/rooms/filter-by-luksuz-and-stdom?st_dom_id=YOUR_STDOM_ID_HERE&luksuzi=klima
```

### 7. Filter by Exact Bed Capacity
```
GET http://localhost:8082/api/v1/rooms/filter-by-krevetnost?exact=2
```

### 8. Filter by Bed Capacity Range
```
GET http://localhost:8082/api/v1/rooms/filter-by-krevetnost?min=2&max=4
```

### 9. Filter by Minimum Bed Capacity
```
GET http://localhost:8082/api/v1/rooms/filter-by-krevetnost?min=3
```

### 10. Search Dormitories by Address
```
GET http://localhost:8082/api/v1/st-doms/search-by-address?address=Studentska
```

### 11. Advanced Filter (All Parameters)
```
GET http://localhost:8082/api/v1/rooms/advanced-filter?luksuzi=klima,terasa&st_dom_id=YOUR_STDOM_ID_HERE&min=2&max=4
```

### 12. Advanced Filter (Luxury + Exact Beds)
```
GET http://localhost:8082/api/v1/rooms/advanced-filter?luksuzi=klima&exact=2
```

---

## Valid Luxury Amenities (luksuzi)

Use these exact values (comma-separated for multiple):
- `klima` - Air conditioning
- `terasa` - Terrace
- `sopstveno kupatilo` - Private bathroom *(URL encode as `sopstveno%20kupatilo`)*
- `áram` - Electric current
- `ablak` - Window
- `neisvrljan zid` - Unpainted wall *(URL encode as `neisvrljan%20zid`)*

---

## How to Get Real IDs

1. Start with: `GET http://localhost:8082/api/v1/st-doms`
2. Copy an `id` value from the response
3. Use that ID in the `st_dom_id` parameter

---

## Postman Setup Steps

1. **Create New Collection** → Name it "Open Data Service"
2. **Add Environment Variable:**
   - Variable: `baseUrl`
   - Value: `http://localhost:8082`
3. **Add 8 Requests** using the URLs above
4. **Replace URLs** with `{{baseUrl}}/api/v1/...` to use the variable

---

## Testing Workflow

```
Step 1: Test health endpoint
  → GET /health

Step 2: Get all data to see what exists
  → GET /api/v1/rooms
  → GET /api/v1/st-doms

Step 3: Copy a real st_dom_id from response

Step 4: Test filtering endpoints with real data
  → Try each filter endpoint
  → Use the copied st_dom_id where needed

Step 5: Test error cases
  → Invalid luxury amenity
  → Invalid ObjectID format
  → Missing required parameters
```

---

## Example Response Format

All successful responses follow this structure:
```json
{
  "data": [...],
  "count": <number>
}
```

Error responses:
```json
{
  "error": "Error message here"
}
```

---

## Status Codes

- `200` - Success
- `400` - Bad Request (invalid parameters)
- `500` - Internal Server Error (database/server issue)

---

## Tips

✅ **Do:**
- Test health check first
- Get all data before filtering
- Use real IDs from your database
- URL encode special characters (spaces, etc.)

❌ **Don't:**
- Try to send JSON body (these are GET requests!)
- Use fake/example IDs (use real ones from your data)
- Forget to URL encode spaces in luxury amenities

