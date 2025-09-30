# Advanced Filter with Address Search

## Overview

The **advanced-filter** endpoint has been enhanced to support **address regex searching**, allowing you to filter rooms not just by specific dormitory ID, but also by address pattern.

## What's New

Previously, you could only filter rooms by a specific `st_dom_id`. Now you can also filter rooms by dormitory address pattern using regex, just like the standalone `search-by-address` endpoint.

## Endpoint

### `GET /api/v1/rooms/advanced-filter`

**Query Parameters:**
- `luksuzi` (string, optional) - Comma-separated list of luxury amenities
- `st_dom_id` (string, optional) - Specific student dormitory ID
- **`address` (string, optional) - Address search pattern (regex, case-insensitive)** ✨ **NEW**
- `exact` (int, optional) - Exact bed capacity
- `min` (int, optional) - Minimum bed capacity
- `max` (int, optional) - Maximum bed capacity

---

## How It Works

When you provide an `address` parameter:
1. **Step 1:** The service searches all dormitories matching the address pattern
2. **Step 2:** Finds all rooms in those matching dormitories
3. **Step 3:** Applies additional filters (luxury, bed capacity) to those rooms
4. **Step 4:** Returns rooms with full dormitory information

**Note:** If you provide both `st_dom_id` and `address`, the `st_dom_id` takes priority (more specific).

---

## Usage Examples

### Example 1: Filter by Address Pattern Only

Find all rooms in dormitories located in "Beograd":

```bash
GET http://localhost:8082/api/v1/rooms/advanced-filter?address=Beograd
```

**Response:**
```json
{
  "data": [
    {
      "id": "room1_id",
      "st_dom_id": "dorm1_id",
      "krevetnost": 2,
      "luksuzi": ["klima"],
      "created_at": "...",
      "updated_at": "...",
      "st_dom": {
        "id": "dorm1_id",
        "ime": "Studentski Dom Karaburma",
        "address": "Studentska 123, Beograd",
        "telephone_number": "+381 11 1234567",
        "email": "karaburma@example.rs",
        "created_at": "...",
        "updated_at": "..."
      }
    },
    {
      "id": "room2_id",
      "st_dom_id": "dorm2_id",
      "krevetnost": 3,
      "luksuzi": ["terasa"],
      "created_at": "...",
      "updated_at": "...",
      "st_dom": {
        "id": "dorm2_id",
        "ime": "Kralj Aleksandar I",
        "address": "Studentski Trg 1, Beograd",
        "telephone_number": "+381 11 9876543",
        "email": "kralj@example.rs",
        "created_at": "...",
        "updated_at": "..."
      }
    }
  ],
  "count": 2
}
```

---

### Example 2: Filter by Address + Luxury Amenities

Find rooms with air conditioning in dormitories on "Studentska" street:

```bash
GET http://localhost:8082/api/v1/rooms/advanced-filter?address=Studentska&luksuzi=klima
```

**Use Case:** "I want a room with A/C in any dorm on Studentska street"

---

### Example 3: Filter by Address + Bed Capacity

Find 2-bed rooms in dormitories in "Novi Sad":

```bash
GET http://localhost:8082/api/v1/rooms/advanced-filter?address=Novi%20Sad&exact=2
```

**Use Case:** "Show me all 2-person rooms in Novi Sad dormitories"

---

### Example 4: Combined Multi-Criteria Search

Find rooms with A/C and terrace, 2-4 beds, in dormitories in "Beograd":

```bash
GET http://localhost:8082/api/v1/rooms/advanced-filter?address=Beograd&luksuzi=klima,terasa&min=2&max=4
```

**Use Case:** "I want a 2-4 person room with A/C and terrace in any Belgrade dormitory"

---

### Example 5: Regex Pattern - Street Names

Find rooms in dormitories on "Bulevar" streets:

```bash
GET http://localhost:8082/api/v1/rooms/advanced-filter?address=^Bulevar
```

**Use Case:** "Show rooms in any dorm on a Boulevard"

---

### Example 6: Regex Pattern - Specific Area

Find rooms in dormitories with "Centar" (downtown) in address:

```bash
GET http://localhost:8082/api/v1/rooms/advanced-filter?address=Centar&luksuzi=sopstveno%20kupatilo
```

**Use Case:** "Rooms with private bathroom in downtown dormitories"

---

## Comparison: Before vs After

### Before (Old Behavior)

You could only filter by **specific dormitory ID**:

```bash
# Only worked if you knew the exact ObjectID
GET /api/v1/rooms/advanced-filter?st_dom_id=507f1f77bcf86cd799439011&luksuzi=klima
```

**Limitation:** You needed to know the exact dormitory ID beforehand.

---

### After (New Behavior)

You can now filter by **address pattern**:

```bash
# Works with any part of the address
GET /api/v1/rooms/advanced-filter?address=Beograd&luksuzi=klima

# Find rooms in dormitories on Studentska street
GET /api/v1/rooms/advanced-filter?address=Studentska&min=2

# Find rooms in Novi Sad dormitories with 2 beds exactly
GET /api/v1/rooms/advanced-filter?address=Novi%20Sad&exact=2
```

**Advantage:** Much more flexible - search by city, street, or any address pattern!

---

## Use Cases

### 1. City-Based Search
Find all rooms in a specific city:
```bash
GET /api/v1/rooms/advanced-filter?address=Beograd
GET /api/v1/rooms/advanced-filter?address=Novi%20Sad
GET /api/v1/rooms/advanced-filter?address=Niš
```

### 2. Street/Area Search
Find rooms in specific streets or areas:
```bash
GET /api/v1/rooms/advanced-filter?address=Studentska
GET /api/v1/rooms/advanced-filter?address=Bulevar%20Oslobodjenja
```

### 3. Combined Location + Amenities
Find rooms with specific amenities in a location:
```bash
GET /api/v1/rooms/advanced-filter?address=Beograd&luksuzi=klima,sopstveno%20kupatilo
```

### 4. Combined Location + Capacity
Find right-sized rooms in a location:
```bash
GET /api/v1/rooms/advanced-filter?address=Centar&min=1&max=2
```

### 5. Full Search
Combine all criteria:
```bash
GET /api/v1/rooms/advanced-filter?address=Novi%20Sad&luksuzi=klima&exact=2
```

---

## Address vs st_dom_id

| Parameter | Type | Use When | Example |
|-----------|------|----------|---------|
| `st_dom_id` | ObjectID | You know exact dormitory | `507f1f77bcf86cd799439011` |
| `address` | Regex Pattern | You want to search by location | `Beograd`, `Studentska`, `^Bulevar` |

**Priority:** If both are provided, `st_dom_id` takes priority (it's more specific).

---

## Performance Notes

### How It Works Internally

1. **With `address` parameter:**
   - Query 1: Search `st_doms` collection by address pattern → Get matching dorm IDs
   - Query 2: Search `sobas` collection for rooms in those dorms + apply other filters
   - Query 3: Fetch full dormitory details for results
   - **Total: ~3 queries**

2. **With `st_dom_id` parameter:**
   - Query 1: Search `sobas` collection for rooms in specific dorm + apply other filters
   - Query 2: Fetch dormitory details
   - **Total: ~2 queries**

**Recommendation:** Use `st_dom_id` when you know it for better performance. Use `address` for flexible searching.

---

## Error Handling

### Empty Results

If no dormitories match the address pattern:
```json
{
  "data": [],
  "count": 0
}
```

### Combined with Other Filters

If address matches dormitories but no rooms match other criteria:
```json
{
  "data": [],
  "count": 0
}
```

---

## cURL Examples

### Basic Address Search
```bash
curl "http://localhost:8082/api/v1/rooms/advanced-filter?address=Beograd"
```

### Address + Amenities
```bash
curl "http://localhost:8082/api/v1/rooms/advanced-filter?address=Beograd&luksuzi=klima,terasa"
```

### Address + Capacity
```bash
curl "http://localhost:8082/api/v1/rooms/advanced-filter?address=Studentska&min=2&max=4"
```

### Full Example
```bash
curl "http://localhost:8082/api/v1/rooms/advanced-filter?address=Novi%20Sad&luksuzi=klima&exact=2"
```

---

## PowerShell Examples

```powershell
# Basic address search
Invoke-RestMethod -Uri "http://localhost:8082/api/v1/rooms/advanced-filter?address=Beograd"

# Address + amenities
Invoke-RestMethod -Uri "http://localhost:8082/api/v1/rooms/advanced-filter?address=Beograd&luksuzi=klima,terasa"

# Full search
Invoke-RestMethod -Uri "http://localhost:8082/api/v1/rooms/advanced-filter?address=Novi Sad&luksuzi=klima&exact=2"
```

---

## Postman Testing

### Request Setup
1. **Method:** `GET`
2. **URL:** `http://localhost:8082/api/v1/rooms/advanced-filter`
3. **Query Params:**
   - `address` = `Beograd`
   - `luksuzi` = `klima` (optional)
   - `min` = `2` (optional)
   - `max` = `4` (optional)

### Test Cases

| Test Case | Parameters | Expected Result |
|-----------|------------|-----------------|
| Address only | `address=Beograd` | All rooms in Belgrade dorms |
| Address + luxury | `address=Beograd&luksuzi=klima` | Rooms with A/C in Belgrade |
| Address + capacity | `address=Novi Sad&exact=2` | 2-bed rooms in Novi Sad |
| Full filter | `address=Beograd&luksuzi=klima&min=2&max=4` | Combined criteria |
| Regex pattern | `address=^Bulevar` | Rooms on boulevards |
| No matches | `address=XYZ123` | Empty array |

---

## Implementation Details

### Service Layer
**File:** `open_data_service/services/open_data_service.go`

The function now:
1. Checks if `addressPattern` is provided
2. If yes, searches dormitories by address regex
3. Extracts matching dormitory IDs
4. Filters rooms by those dormitory IDs (+ other criteria)
5. Returns rooms with full dormitory info

### Handler Layer
**File:** `open_data_service/handlers/open_data_handler.go`

Added:
```go
// Parse address pattern
addressPattern := c.Query("address")
```

Updated service call:
```go
rooms, err := h.service.AdvancedFilterRooms(luksuzi, stDomID, addressPattern, exact, min, max)
```

---

## Files Modified

1. ✅ `open_data_service/services/open_data_service.go` - Updated `AdvancedFilterRooms()` method
2. ✅ `open_data_service/handlers/open_data_handler.go` - Updated `AdvancedFilterRooms()` handler

---

## Summary

The **advanced-filter** endpoint is now even more powerful:

✅ Filter rooms by luxury amenities (luksuzi)  
✅ Filter by specific dormitory ID (st_dom_id)  
✅ **Filter by address pattern (address)** ← **NEW!**  
✅ Filter by bed capacity (krevetnost)  
✅ Combine any/all criteria  
✅ Case-insensitive regex support  
✅ Returns full dormitory information

Perfect for flexible location-based room searches! 🎯

