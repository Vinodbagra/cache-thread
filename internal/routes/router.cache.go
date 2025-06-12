package routes

import (
	"time"

	"github.com/Vinodbagra/cache-thread/internal/handler"
	"github.com/Vinodbagra/cache-thread/internal/service"
	"github.com/gin-gonic/gin"
)

type cacheRoutes struct {
	Handler *handler.CacheHandler
	router  *gin.RouterGroup
}

func NewCacheRoute(router *gin.RouterGroup, cacheMaxSize int, cacheDefaultTTL time.Duration) *cacheRoutes {
	cacheService := service.NewCacheService(cacheMaxSize, cacheDefaultTTL)
	cacheHandler := handler.NewCacheHandler(cacheService)

	return &cacheRoutes{Handler: cacheHandler, router: router}
}

func (r *cacheRoutes) Routes() {
	// Cache API Routes
	cacheRoute := r.router.Group("/cache")
	{
		// Basic CRUD operations
		cacheRoute.PUT("/put", r.Handler.Put)               // Store key-value pair
		cacheRoute.GET("/get/:key", r.Handler.Get)          // Get value by key
		cacheRoute.DELETE("/delete/:key", r.Handler.Delete) // Delete key
		cacheRoute.DELETE("/clear", r.Handler.Clear)        // Clear entire cache

		// Bulk operations
		cacheRoute.POST("/bulk/put", r.Handler.BulkPut) // Bulk store key-value pairs
		cacheRoute.POST("/bulk/get", r.Handler.BulkGet) // Bulk get values

		// Information and monitoring
		cacheRoute.GET("/stats", r.Handler.GetStats)          // Get cache statistics
		cacheRoute.GET("/health", r.Handler.GetHealth)        // Health check
		cacheRoute.GET("/keys", r.Handler.GetKeys)            // List all keys (for debugging)
		cacheRoute.GET("/config", r.Handler.GetConfiguration) // Get cache configuration
	}
}
