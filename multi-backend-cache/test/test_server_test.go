package test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"multi-backend-cache/Internal/cache"
// 	"multi-backend-cache/Internal/server"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// // setupRouters initializes the server with all routes
// func setupRouters() *gin.Engine {
// 	gin.SetMode(gin.TestMode)
// 	// Initialize mock caches with hypothetical functions
// 	inmemoryCache := cache.NewLRUCache(100, 60) // Assuming this function exists and works
// 	redisCache := cache.NewRedisCache("localhost:6379", "", 0, 1*time.Minute)
// 	memCache := cache.NewMemCache("localhost:11211", 60)

// 	cacheSystemType := server.NewServer(inmemoryCache, redisCache, memCache)
// 	router := gin.Default()

// 	// Setup routes
// 	router.GET("/cache/:key", cacheSystemType.GetCacheHandler)
// 	router.POST("/cache", cacheSystemType.SetCacheHandler)
// 	router.DELETE("/cache/:key", cacheSystemType.DeleteCacheHandler)

// 	return router
// }

// // testSetAndGetCacheHandler tests the SetCacheHandler and GetCacheHandler in a single function.
// func testSetAndGetCacheHandler(t *testing.T, key, value string, ttl int) {
// 	router := setupRouters()

// 	// Define the payload for setting a value
// 	payload := map[string]interface{}{
// 		"key":   key,
// 		"value": value,
// 		"ttl":   ttl,
// 	}
// 	body, err := json.Marshal(payload)
// 	assert.Nil(t, err)

// 	// Send a POST request to set the cache value
// 	req, _ := http.NewRequest("POST", "/cache?cache=redis", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// Verify the response for setting the cache
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	var response map[string]string
// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "ok", response["status"])

// 	// Send a GET request to retrieve the cache value
// 	req, _ = http.NewRequest("GET", fmt.Sprintf("/cache/%s?cache=redis", key), nil)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// Verify the response for getting the cache
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Attempt to unmarshal the response as a JSON object
// 	var getResponse map[string]interface{}
// 	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
// 	if err != nil {
// 		// If unmarshalling fails, try to interpret the response as a direct string
// 		var directResponse string
// 		err = json.Unmarshal(w.Body.Bytes(), &directResponse)
// 		assert.Nil(t, err)
// 		assert.Equal(t, value, directResponse)
// 	} else {
// 		// If unmarshalling as JSON succeeds, verify the value
// 		assert.Equal(t, value, getResponse["value"])
// 	}
// }

// // testDeleteCacheHandler tests the DeleteCacheHandler with a single request.
// func testDeleteCacheHandler(t *testing.T, key, value string, ttl int) {
// 	router := setupRouters()

// 	// Set a value first
// 	payload := map[string]interface{}{
// 		"key":   key,
// 		"value": value,
// 		"ttl":   ttl,
// 	}
// 	body, err := json.Marshal(payload)
// 	assert.Nil(t, err)

// 	req, _ := http.NewRequest("POST", "/cache?cache=redis", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Now delete the value
// 	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/cache/%s?cache=redis", key), nil)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]string
// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "ok", response["status"])

// 	// Sending another GET request to validate if there are any non-deleted values.
// 	req, _ = http.NewRequest("GET", fmt.Sprintf("/cache/%s?cache=redis", key), nil)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// }

// func TestSetGetAndDeleteCache(t *testing.T) {
// 	key := "Kannan"
// 	value := "Software Engineer"
// 	ttl := 60
// 	testSetAndGetCacheHandler(t, key, value, ttl)
// 	testDeleteCacheHandler(t, key, value, ttl)
// }
