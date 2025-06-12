# Cache API Endpoints

## Environment Variables

Create a `.env` file with the following variables:

```env
# Server Configuration
PORT=8080
ENVIRONMENT=development
DEBUG=true

# Cache Configuration
CACHE_MAX_SIZE=1000
CACHE_TTL=30m

# Database Configuration (optional - only if using database)
DB_POSTGRE_DSN=postgres://username:password@localhost:5432/cache_db?sslmode=disable
DB_POSTGRE_URL=postgres://username:password@localhost:5432/cache_db?sslmode=disable
```

## API Endpoints

Base URL: `http://localhost:8080/api/cache`

### Basic CRUD Operations

#### 1. Store Key-Value Pair
- **Method:** `PUT`
- **Endpoint:** `/put`
- **Body:**
```json
{
  "key": "user:123",
  "value": {"name": "John Doe", "email": "john@example.com"},
  "ttl": 3600
}
```

#### 2. Get Value by Key
- **Method:** `GET`
- **Endpoint:** `/get/{key}`
- **Example:** `/get/user:123`

#### 3. Delete Key
- **Method:** `DELETE`
- **Endpoint:** `/delete/{key}`
- **Example:** `/delete/user:123`

#### 4. Clear Entire Cache
- **Method:** `DELETE`
- **Endpoint:** `/clear`

### Bulk Operations

#### 5. Bulk Store Key-Value Pairs
- **Method:** `POST`
- **Endpoint:** `/bulk/put`
- **Body:**
```json
{
  "items": [
    {
      "key": "user:1",
      "value": {"name": "Alice"},
      "ttl": 1800
    },
    {
      "key": "user:2", 
      "value": {"name": "Bob"},
      "ttl": 3600
    }
  ]
}
```

#### 6. Bulk Get Values
- **Method:** `POST`
- **Endpoint:** `/bulk/get`
- **Body:**
```json
{
  "keys": ["user:1", "user:2", "user:3"]
}
```

### Information and Monitoring

#### 7. Get Cache Statistics
- **Method:** `GET`
- **Endpoint:** `/stats`
- **Response:**
```json
{
  "hits": 150,
  "misses": 25,
  "hit_rate": 0.857,
  "total_requests": 175,
  "current_size": 45,
  "max_size": 1000,
  "evictions": 5,
  "expired_removals": 10,
  "uptime": "2h30m15s"
}
```

#### 8. Health Check
- **Method:** `GET`
- **Endpoint:** `/health`
- **Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "uptime": "2h30m15s"
}
```

#### 9. List All Keys (Debug)
- **Method:** `GET`
- **Endpoint:** `/keys`
- **Query Parameters:**
  - `limit` (optional): Maximum number of keys to return (default: 100)
- **Example:** `/keys?limit=50`

#### 10. Get Cache Configuration
- **Method:** `GET`
- **Endpoint:** `/config`
- **Response:**
```json
{
  "max_size": 1000,
  "default_ttl": "30m0s",
  "cleanup_interval": "30s",
  "start_time": "2024-01-15T08:00:00Z",
  "uptime": "2h30m15s"
}
```

## Response Formats

### Success Responses
- **200 OK:** Operation completed successfully
- **201 Created:** Resource created successfully (for PUT operations)

### Error Responses
```json
{
  "error": "Error description",
  "code": "ERROR_CODE",
  "message": "Detailed error message"
}
```

### Common Error Codes
- `INVALID_REQUEST`: Invalid request body or parameters
- `MISSING_KEY`: Key parameter is missing or empty
- `PUT_FAILED`: Failed to store key-value pair
- `EMPTY_REQUEST`: No items or keys provided in bulk operations

## Features

- **LRU Eviction:** Least Recently Used items are evicted when cache is full
- **TTL Support:** Automatic expiration of cached items
- **Bulk Operations:** Efficient batch processing
- **Statistics:** Real-time cache performance metrics
- **Thread-Safe:** Concurrent access support
- **Background Cleanup:** Automatic removal of expired items 