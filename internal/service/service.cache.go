package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/Vinodbagra/cache-thread/internal/models"
)

// CacheService implements the cache business logic
type CacheService struct {
	data         map[string]*models.CacheEntry
	head         *models.CacheEntry // Most recently used
	tail         *models.CacheEntry // Least recently used
	maxSize      int
	defaultTTL   time.Duration
	startTime    time.Time
	
	// Statistics
	hits            int64
	misses          int64
	evictions       int64
	expiredRemovals int64
	
	// Synchronization
	mutex       sync.RWMutex
	cleanupDone chan bool
	stopCleanup chan bool
}

// NewCacheService creates a new cache service instance
func NewCacheService(maxSize int, defaultTTL time.Duration) *CacheService {
	service := &CacheService{
		data:        make(map[string]*models.CacheEntry),
		maxSize:     maxSize,
		defaultTTL:  defaultTTL,
		startTime:   time.Now(),
		cleanupDone: make(chan bool),
		stopCleanup: make(chan bool),
	}
	
	// Initialize doubly linked list with sentinel nodes
	service.head = &models.CacheEntry{}
	service.tail = &models.CacheEntry{}
	service.head.Next = service.tail
	service.tail.Prev = service.head
	
	// Start background cleanup goroutine
	go service.cleanupWorker()
	
	return service
}

// Put inserts or updates a key-value pair with optional TTL
func (cs *CacheService) Put(key string, value interface{}, ttl *time.Duration) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	
	var expiration int64
	if ttl != nil && *ttl > 0 {
		expiration = time.Now().Add(*ttl).Unix()
	} else if cs.defaultTTL > 0 {
		expiration = time.Now().Add(cs.defaultTTL).Unix()
	}
	
	now := time.Now()
	
	if entry, exists := cs.data[key]; exists {
		// Update existing entry
		entry.Value = value
		entry.Expiration = expiration
		entry.AccessedAt = now
		cs.moveToHead(entry)
	} else {
		// Create new entry
		entry := &models.CacheEntry{
			Key:        key,
			Value:      value,
			Expiration: expiration,
			CreatedAt:  now,
			AccessedAt: now,
		}
		
		// Check if we need to evict
		if len(cs.data) >= cs.maxSize {
			cs.evictLRU()
		}
		
		cs.data[key] = entry
		cs.addToHead(entry)
	}
	
	return nil
}

// Get retrieves a value by key and updates access order
func (cs *CacheService) Get(key string) (*models.CacheEntry, bool) {
	if key == "" {
		return nil, false
	}
	
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	
	entry, exists := cs.data[key]
	if !exists {
		cs.misses++
		return nil, false
	}
	
	// Check if entry has expired
	if entry.IsExpired() {
		cs.removeEntry(entry)
		cs.expiredRemovals++
		cs.misses++
		return nil, false
	}
	
	// Update access time and move to head (most recently used)
	entry.UpdateAccessTime()
	cs.moveToHead(entry)
	cs.hits++
	
	return entry, true
}

// Delete removes a specific key from the cache
func (cs *CacheService) Delete(key string) (bool, bool) {
	if key == "" {
		return false, false
	}
	
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	
	entry, exists := cs.data[key]
	if !exists {
		return false, false
	}
	
	cs.removeEntry(entry)
	return true, true
}

// Clear removes all entries from the cache
func (cs *CacheService) Clear() int {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	
	itemsCleared := len(cs.data)
	cs.data = make(map[string]*models.CacheEntry)
	cs.head.Next = cs.tail
	cs.tail.Prev = cs.head
	
	return itemsCleared
}

// GetStats returns current cache statistics
func (cs *CacheService) GetStats() models.CacheStats {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	
	totalRequests := cs.hits + cs.misses
	var hitRate float64
	if totalRequests > 0 {
		hitRate = float64(cs.hits) / float64(totalRequests)
	}
	
	uptime := time.Since(cs.startTime).String()
	
	return models.CacheStats{
		Hits:            cs.hits,
		Misses:          cs.misses,
		HitRate:         hitRate,
		TotalRequests:   totalRequests,
		CurrentSize:     len(cs.data),
		MaxSize:         cs.maxSize,
		Evictions:       cs.evictions,
		ExpiredRemovals: cs.expiredRemovals,
		Uptime:          uptime,
	}
}

// GetConfiguration returns cache configuration
func (cs *CacheService) GetConfiguration() models.CacheConfiguration {
	return models.CacheConfiguration{
		MaxSize:         cs.maxSize,
		DefaultTTL:      cs.defaultTTL,
		CleanupInterval: 30 * time.Second,
		StartTime:       cs.startTime,
	}
}

// BulkPut performs multiple put operations
func (cs *CacheService) BulkPut(items []models.PutRequest) models.BulkPutResponse {
	response := models.BulkPutResponse{}
	
	for _, item := range items {
		var ttl *time.Duration
		if item.TTL != nil && *item.TTL > 0 {
			duration := time.Duration(*item.TTL) * time.Second
			ttl = &duration
		}
		
		if err := cs.Put(item.Key, item.Value, ttl); err != nil {
			response.Failed++
			response.Errors = append(response.Errors, fmt.Sprintf("Key '%s': %v", item.Key, err))
		} else {
			response.Successful++
		}
	}
	
	return response
}

// BulkGet performs multiple get operations
func (cs *CacheService) BulkGet(keys []string) models.BulkGetResponse {
	response := models.BulkGetResponse{
		Results: make(map[string]models.GetResponse),
	}
	
	for _, key := range keys {
		if entry, found := cs.Get(key); found {
			response.Results[key] = entry.ToResponse()
			response.Found++
		} else {
			response.Results[key] = models.GetResponse{
				Key:   key,
				Found: false,
			}
			response.NotFound++
		}
	}
	
	return response
}

// ListKeys returns all keys in the cache (for debugging)
func (cs *CacheService) ListKeys() []string {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	
	keys := make([]string, 0, len(cs.data))
	for key := range cs.data {
		keys = append(keys, key)
	}
	
	return keys
}

// Close stops the background cleanup worker
func (cs *CacheService) Close() {
	close(cs.stopCleanup)
	<-cs.cleanupDone
}

// Internal methods for LRU management

// addToHead adds a new entry right after head (most recently used position)
func (cs *CacheService) addToHead(entry *models.CacheEntry) {
	entry.Prev = cs.head
	entry.Next = cs.head.Next
	cs.head.Next.Prev = entry
	cs.head.Next = entry
}

// removeFromList removes an entry from the doubly linked list
func (cs *CacheService) removeFromList(entry *models.CacheEntry) {
	entry.Prev.Next = entry.Next
	entry.Next.Prev = entry.Prev
}

// moveToHead moves an existing entry to head (mark as most recently used)
func (cs *CacheService) moveToHead(entry *models.CacheEntry) {
	cs.removeFromList(entry)
	cs.addToHead(entry)
}

// evictLRU removes the least recently used entry
func (cs *CacheService) evictLRU() {
	if cs.tail.Prev != cs.head {
		lru := cs.tail.Prev
		cs.removeEntry(lru)
		cs.evictions++
	}
}

// removeEntry removes an entry from both map and linked list
func (cs *CacheService) removeEntry(entry *models.CacheEntry) {
	delete(cs.data, entry.Key)
	cs.removeFromList(entry)
}

// cleanupWorker runs periodically to remove expired entries
func (cs *CacheService) cleanupWorker() {
	ticker := time.NewTicker(30 * time.Second) // Cleanup every 30 seconds
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			cs.cleanupExpired()
		case <-cs.stopCleanup:
			cs.cleanupDone <- true
			return
		}
	}
}

// cleanupExpired removes all expired entries
func (cs *CacheService) cleanupExpired() {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	
	var expiredKeys []string
	for key, entry := range cs.data {
		if entry.IsExpired() {
			expiredKeys = append(expiredKeys, key)
		}
	}
	
	for _, key := range expiredKeys {
		if entry, exists := cs.data[key]; exists {
			cs.removeEntry(entry)
			cs.expiredRemovals++
		}
	}
}