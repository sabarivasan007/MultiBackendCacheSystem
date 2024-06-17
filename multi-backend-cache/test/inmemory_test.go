package test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"multi-backend-cache/Internal/cache"
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// func TestIsExpired(t *testing.T) {
// 	a, _ := time.Parse(time.RFC3339Nano, "2024-06-10T13:38:18.616332+05:30")
// 	got := cache.IsExpired(a)
// 	want := true
// 	assert.Equal(t, want, got)
// }

// // func TestCalculateExpiryTime(t *testing.T) {
// // 	inmemoryCache := cache.NewLRUCache(300, 60)
// // 	got := CalculateExpiryTime(60)
// // 	want := time.Now().Add(60 * time.Second)
// // 	assert.Equal(t, want, got)
// // }

// func TestCalculateSize(t *testing.T) {
// 	// Create a CacheData instance
// 	cacheData := &cache.CacheData{
// 		Key:        "test_key",
// 		Value:      "test_value",
// 		TTL:        60,
// 		ExpiryTime: time.Now().Add(time.Second * 10),
// 	}

// 	// Call the calculateSize function
// 	size := cache.CalculateSize(cacheData)
// 	fmt.Println("size", size)

// 	// Expected size
// 	expectedSize := int(reflect.TypeOf(*cacheData).Size())
// 	fmt.Println("expected_ size", expectedSize)

// 	// Check if the size is as expected
// 	if size != expectedSize {
// 		t.Errorf("Expected size: %d, got: %d", expectedSize, size)
// 	}
// }

// func setupRouter() *gin.Engine {
// 	inmemoryCache := cache.NewLRUCache(300, 60)
// 	cacheSystemType := server.NewServer(inmemoryCache, nil, nil)
// 	router := gin.Default()

// 	router.GET("/getCache/:key", cacheSystemType.GetCacheHandler)
// 	router.POST("/setCache", cacheSystemType.SetCacheHandler)
// 	router.DELETE("/deleteCache/:key", cacheSystemType.DeleteCacheHandler)

// 	return router
// }

// func TestPostCacheHandler(t *testing.T) {
// 	router := setupRouter()

// 	w := httptest.NewRecorder()
// 	reqBody := `{"key": "1", "value": {"id":"12345","name":"Abcd"}, "ttl": 60}`
// 	req, _ := http.NewRequest("POST", "/setCache", strings.NewReader(reqBody))
// 	fmt.Println("req", req)
// 	req.Header.Set("Content-Type", "application/json")
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Parse the response body
// 	var response cache.CacheData
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)

// 	// Calculate the expected expiry time
// 	expectedExpiryTime := time.Now().Add(60 * time.Second)

// 	// Allow a margin of error in the expiry time due to processing delays
// 	marginOfError := 1 * time.Second
// 	assert.WithinDuration(t, expectedExpiryTime, response.ExpiryTime, marginOfError)

// 	assert.Equal(t, "1", response.Key)
// 	expectedValues := map[string]interface{}{
// 		"id":   "12345",
// 		"name": "Abcd",
// 	}
// 	fmt.Println("expected", expectedValues)
// 	fmt.Println("actual val", response.Value)
// 	assert.Equal(t, expectedValues, response.Value)
// }

// func TestGetCacheHandler(t *testing.T) {
// 	router := setupRouter()

// 	// First, post a cache entry
// 	w := httptest.NewRecorder()
// 	reqBody := `{"key": "2", "value": "session", "ttl": 300}`
// 	req, _ := http.NewRequest("POST", "/setCache", strings.NewReader(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Then, get the cache entry
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("GET", "/getCache/2", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Parse the response body
// 	var response cache.CacheData
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)

// 	// Calculate the expected expiry time
// 	expectedExpiryTime := time.Now().Add(300 * time.Second)

// 	// Allow a margin of error in the expiry time due to processing delays
// 	marginOfError := 2 * time.Second
// 	assert.WithinDuration(t, expectedExpiryTime, response.ExpiryTime, marginOfError)

// 	assert.Equal(t, "2", response.Key)
// 	assert.Equal(t, "session", response.Value)
// }

// func TestDeleteCacheHandler(t *testing.T) {
// 	router := setupRouter()

// 	// First, post a cache entry
// 	w := httptest.NewRecorder()
// 	reqBody := `{"key": "3", "value": "cache", "ttl": 300}`
// 	req, _ := http.NewRequest("POST", "/setCache", strings.NewReader(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Then, delete the cache entry
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("DELETE", "/deleteCache/3", nil)
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), `"message":"Cache 3 is successfully deleted"`)

// 	// Ensure the cache entry is deleted
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("GET", "/getCache/3", nil)
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusNotFound, w.Code)
// }
