package models

import "time"

// CacheEntry represents a single cache entry with value, expiration time, and LRU pointers
type CacheEntry struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Expiration int64       `json:"expiration"` // Unix timestamp, 0 means no expiration
	CreatedAt  time.Time   `json:"created_at"`
	AccessedAt time.Time   `json:"accessed_at"`
	Prev       *CacheEntry
	Next       *CacheEntry
}

// CacheStats holds statistics about cache performance
type CacheStats struct {
	Hits            int64   `json:"hits"`
	Misses          int64   `json:"misses"`
	HitRate         float64 `json:"hit_rate"`
	TotalRequests   int64   `json:"total_requests"`
	CurrentSize     int     `json:"current_size"`
	MaxSize         int     `json:"max_size"`
	Evictions       int64   `json:"evictions"`
	ExpiredRemovals int64   `json:"expired_removals"`
	Uptime          string  `json:"uptime"`
}

// PutRequest represents the request body for PUT operations
type PutRequest struct {
	Key   string      `json:"key" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
	TTL   *int        `json:"ttl,omitempty"` // TTL in seconds, optional
}

// GetResponse represents the response for GET operations
type GetResponse struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Found      bool        `json:"found"`
	Expired    bool        `json:"expired,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
	AccessedAt time.Time   `json:"accessed_at,omitempty"`
}

// DeleteResponse represents the response for DELETE operations
type DeleteResponse struct {
	Key     string `json:"key"`
	Deleted bool   `json:"deleted"`
	Found   bool   `json:"found"`
}

// ClearResponse represents the response for CLEAR operations
type ClearResponse struct {
	ItemsCleared int    `json:"items_cleared"`
	Message      string `json:"message"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

// BulkPutRequest represents bulk put operations
type BulkPutRequest struct {
	Items []PutRequest `json:"items" binding:"required"`
}

// BulkPutResponse represents bulk put response
type BulkPutResponse struct {
	Successful int      `json:"successful"`
	Failed     int      `json:"failed"`
	Errors     []string `json:"errors,omitempty"`
}

// BulkGetRequest represents bulk get operations
type BulkGetRequest struct {
	Keys []string `json:"keys" binding:"required"`
}

// BulkGetResponse represents bulk get response
type BulkGetResponse struct {
	Results map[string]GetResponse `json:"results"`
	Found   int                    `json:"found"`
	NotFound int                   `json:"not_found"`
}

// CacheConfiguration represents cache configuration
type CacheConfiguration struct {
	MaxSize         int           `json:"max_size"`
	DefaultTTL      time.Duration `json:"default_ttl"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
	StartTime       time.Time     `json:"start_time"`
}

// IsExpired checks if the cache entry has expired
func (ce *CacheEntry) IsExpired() bool {
	if ce.Expiration == 0 {
		return false // No expiration set
	}
	return time.Now().Unix() > ce.Expiration
}

// UpdateAccessTime updates the last accessed time
func (ce *CacheEntry) UpdateAccessTime() {
	ce.AccessedAt = time.Now()
}

// SetExpiration sets the expiration time
func (ce *CacheEntry) SetExpiration(ttl time.Duration) {
	if ttl > 0 {
		ce.Expiration = time.Now().Add(ttl).Unix()
	} else {
		ce.Expiration = 0
	}
}

// GetTTL returns the remaining TTL in seconds
func (ce *CacheEntry) GetTTL() int64 {
	if ce.Expiration == 0 {
		return -1 // No expiration
	}
	remaining := ce.Expiration - time.Now().Unix()
	if remaining < 0 {
		return 0 // Expired
	}
	return remaining
}

// ToResponse converts CacheEntry to GetResponse
func (ce *CacheEntry) ToResponse() GetResponse {
	return GetResponse{
		Key:        ce.Key,
		Value:      ce.Value,
		Found:      true,
		Expired:    ce.IsExpired(),
		CreatedAt:  ce.CreatedAt,
		AccessedAt: ce.AccessedAt,
	}
}