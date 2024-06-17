package main

import (
	"fmt"
	"log"
	handler "multi-backend-cache/Internal/Handler"
	"multi-backend-cache/Internal/cache"
	"multi-backend-cache/Internal/config"
	"multi-backend-cache/Internal/metrices"
	_ "multi-backend-cache/docs"
	"time"

	"github.com/pbnjay/memory"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var inmemorycache *cache.LRUCache
var tenantCaches *cache.FixedTenantsCaches

func DefaultModelsExpandDepth(depth int) func(*ginSwagger.Config) {
	return func(c *ginSwagger.Config) {
		c.DefaultModelsExpandDepth = depth
	}
}

func main() {
	config.LoadConfig("./Internal/config/config.yaml")
	defaultTTL := viper.GetInt("defaultTTL")
	redisCache := cache.NewRedisCache("redis:6379", "", 0, 1*time.Minute)
	memCache := cache.NewMemCache("memcached:11211", 60)
	totalCacheMemory := int(float64(memory.TotalMemory()) * viper.GetFloat64("MemoryUsagePercentage"))

	fmt.Println("totalCacheMemory", totalCacheMemory)
	// Initialize tenant-specific in-memory LRUCaches
	// Each tenant has a cache with a fixed capacity
	isTenantBased := config.AppConfig.IsTenantBased
	if isTenantBased {
		tenantCaches = cache.NewFixedTenantsCaches(totalCacheMemory, time.Duration(defaultTTL))
	} else {
		// inmemorycache = cache.NewLRUCache(300, 60)
		inmemorycache = cache.NewLRUCache(totalCacheMemory, time.Duration(defaultTTL))
	}
	cacheSystem := handler.NewServer(tenantCaches, redisCache, memCache, inmemorycache)

	// cacheSystemType := server.NewServer(tenantCaches, redisCache, memCache)

	router := gin.Default()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
	//     ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
	//     ginSwagger.DefaultModelsExpandDepth(-1), // Set to -1 to hide models completely
	// ),
	// )

	metrics := metrices.NewMetrics()

	router.Use(metrics.Middleware())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(handler.ValidateCacheSystem())

	// Middleware for "inmemory" system
	router.Use(func(c *gin.Context) {
		if system := c.Query("system"); system == "inmemory" && isTenantBased {
			handler.ValidateTenant()(c)
		}
		c.Next()
	})

	// Cache System routes
	router.GET("/cache/:key", cacheSystem.GetCacheHandler)
	//router.GET("/cache/TTL/:key", cacheSystem.GetCacheWithTTLHandler)
	router.POST("/cache", cacheSystem.SetCacheHandler)
	router.DELETE("/cache/:key", cacheSystem.DeleteCacheHandler)
	router.PUT("/cache/clear", cacheSystem.ClearCacheHandler)

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server started at %s\n", addr)
	log.Fatal(router.Run(addr))
}
