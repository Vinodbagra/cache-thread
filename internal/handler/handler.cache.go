package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Vinodbagra/cache-thread/internal/models"
	"github.com/Vinodbagra/cache-thread/internal/service"
	"github.com/gin-gonic/gin"
)



type CacheHandler struct {
	cacheService *service.CacheService
}

func NewCacheHandler(cacheService *service.CacheService) *CacheHandler {
	return &CacheHandler{cacheService: cacheService}
}

func (ch *CacheHandler) Put(c *gin.Context) {
	var req models.PutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	var ttl *time.Duration
	if req.TTL != nil && *req.TTL > 0 {
		duration := time.Duration(*req.TTL) * time.Second
		ttl = &duration
	}

	if err := ch.cacheService.Put(req.Key, req.Value, ttl); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to store key-value pair",
			Code:    "PUT_FAILED",
			Message: err.Error(),
		})
		return
	}

	response := gin.H{
		"message": "Key-value pair stored successfully",
		"key":     req.Key,
		"ttl":     req.TTL,
	}

	c.JSON(http.StatusCreated, response)
}

// Get handles GET requests to retrieve values by key
// @Summary Get value by key
// @Description Retrieve a value from cache by key
// @Tags cache
// @Produce json
// @Param key path string true "Cache key"
// @Success 200 {object} models.GetResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/cache/get/{key} [get]
func (ch *CacheHandler) Get(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Key parameter is required",
			Code:    "MISSING_KEY",
			Message: "Please provide a valid key parameter",
		})
		return
	}

	entry, found := ch.cacheService.Get(key)
	if !found {
		c.JSON(http.StatusNotFound, models.GetResponse{
			Key:   key,
			Found: false,
		})
		return
	}

	response := entry.ToResponse()
	c.JSON(http.StatusOK, response)
}

// Delete handles DELETE requests to remove keys
// @Summary Delete key from cache
// @Description Remove a key-value pair from cache
// @Tags cache
// @Produce json
// @Param key path string true "Cache key"
// @Success 200 {object} models.DeleteResponse
// @Failure 404 {object} models.DeleteResponse
// @Router /api/v1/cache/delete/{key} [delete]
func (ch *CacheHandler) Delete(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Key parameter is required",
			Code:    "MISSING_KEY",
			Message: "Please provide a valid key parameter",
		})
		return
	}

	deleted, found := ch.cacheService.Delete(key)
	
	response := models.DeleteResponse{
		Key:     key,
		Deleted: deleted,
		Found:   found,
	}

	if found {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusNotFound, response)
	}
}

// Clear handles DELETE requests to clear entire cache
// @Summary Clear entire cache
// @Description Remove all key-value pairs from cache
// @Tags cache
// @Produce json
// @Success 200 {object} models.ClearResponse
// @Router /api/v1/cache/clear [delete]
func (ch *CacheHandler) Clear(c *gin.Context) {
	itemsCleared := ch.cacheService.Clear()
	
	response := models.ClearResponse{
		ItemsCleared: itemsCleared,
		Message:      "Cache cleared successfully",
	}

	c.JSON(http.StatusOK, response)
}

// GetStats handles GET requests for cache statistics
// @Summary Get cache statistics
// @Description Retrieve current cache performance statistics
// @Tags cache
// @Produce json
// @Success 200 {object} models.CacheStats
// @Router /api/v1/cache/stats [get]
func (ch *CacheHandler) GetStats(c *gin.Context) {
	stats := ch.cacheService.GetStats()
	c.JSON(http.StatusOK, stats)
}

// BulkPut handles bulk PUT operations
// @Summary Bulk store key-value pairs
// @Description Store multiple key-value pairs in a single request
// @Tags cache
// @Accept json
// @Produce json
// @Param request body models.BulkPutRequest true "Bulk put request"
// @Success 200 {object} models.BulkPutResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/cache/bulk/put [post]
func (ch *CacheHandler) BulkPut(c *gin.Context) {
	var req models.BulkPutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "No items provided",
			Code:    "EMPTY_REQUEST",
			Message: "At least one item must be provided",
		})
		return
	}

	response := ch.cacheService.BulkPut(req.Items)
	c.JSON(http.StatusOK, response)
}

// BulkGet handles bulk GET operations
// @Summary Bulk get values by keys
// @Description Retrieve multiple values from cache by keys
// @Tags cache
// @Accept json
// @Produce json
// @Param request body models.BulkGetRequest true "Bulk get request"
// @Success 200 {object} models.BulkGetResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/cache/bulk/get [post]
func (ch *CacheHandler) BulkGet(c *gin.Context) {
	var req models.BulkGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if len(req.Keys) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "No keys provided",
			Code:    "EMPTY_REQUEST",
			Message: "At least one key must be provided",
		})
		return
	}

	response := ch.cacheService.BulkGet(req.Keys)
	c.JSON(http.StatusOK, response)
}

// GetHealth handles health check requests
// @Summary Health check
// @Description Check if the cache service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /api/v1/health [get]
func (ch *CacheHandler) GetHealth(c *gin.Context) {
	config := ch.cacheService.GetConfiguration()
	
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(config.StartTime).String(),
	}

	c.JSON(http.StatusOK, response)
}

// GetKeys handles requests to list all keys (for debugging)
// @Summary List all keys
// @Description Get list of all keys in cache (for debugging purposes)
// @Tags cache
// @Produce json
// @Param limit query int false "Limit number of keys returned" default(100)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cache/keys [get]
func (ch *CacheHandler) GetKeys(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	allKeys := ch.cacheService.ListKeys()
	
	// Apply limit
	if len(allKeys) > limit {
		allKeys = allKeys[:limit]
	}

	response := gin.H{
		"keys":       allKeys,
		"count":      len(allKeys),
		"limited":    len(ch.cacheService.ListKeys()) > limit,
		"total_keys": len(ch.cacheService.ListKeys()),
	}

	c.JSON(http.StatusOK, response)
}

// GetConfiguration handles requests for cache configuration
// @Summary Get cache configuration
// @Description Retrieve current cache configuration settings
// @Tags cache
// @Produce json
// @Success 200 {object} models.CacheConfiguration
// @Router /api/v1/cache/config [get]
func (ch *CacheHandler) GetConfiguration(c *gin.Context) {
	config := ch.cacheService.GetConfiguration()
	
	// Convert to a more readable format
	response := gin.H{
		"max_size":         config.MaxSize,
		"default_ttl":      config.DefaultTTL.String(),
		"cleanup_interval": config.CleanupInterval.String(),
		"start_time":       config.StartTime,
		"uptime":           time.Since(config.StartTime).String(),
	}

	c.JSON(http.StatusOK, response)
}

