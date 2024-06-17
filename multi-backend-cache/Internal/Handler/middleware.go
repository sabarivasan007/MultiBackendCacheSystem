package handler

import (
	"multi-backend-cache/Internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("tenantID")
		if tenantID == "" {
			tenantID = c.GetHeader("X-Tenant-ID")
		}

		if tenantID == "" || !isTenantValid(tenantID) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tenant Not Found"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isTenantValid(tenantID string) bool {
	for _, tenant := range config.AppConfig.TenantIDs {
		if tenant == tenantID {
			return true
		}
	}
	return false
}

func ValidateCacheSystem() gin.HandlerFunc {
	return func(c *gin.Context) {
		cacheSystem := c.Query("system")
		if cacheSystem == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cache system not provided"})
			c.Abort()
			return
		}

		cacheSystems := config.AppConfig.CacheSystems
		found := false
		for _, sys := range cacheSystems {
			if sys == cacheSystem {
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cache system"})
			c.Abort()
			return
		}

		c.Next()
	}
}
