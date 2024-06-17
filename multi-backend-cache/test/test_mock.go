package test

// import (
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"multi-backend-cache/Internal/cache"
// 	"multi-backend-cache/Internal/server"

// 	"github.com/golang/mock/gomock"
// 	"github.com/gorilla/mux"
// )

// func setupRouter(redisCache cache.CacheSystem, memCache cache.CacheSystem) *mux.Router {
// 	srv := server.NewServer(redisCache, memCache)
// 	r := mux.NewRouter()

// 	r.HandleFunc("/cache/{key}", srv.GetCache).Methods("GET")
// 	r.HandleFunc("/cache/TTL/{key}", srv.GetCacheWithTTL).Methods("GET")
// 	r.HandleFunc("/cache", srv.SetCache).Methods("POST")
// 	r.HandleFunc("/cache/{key}", srv.DeleteCache).Methods("DELETE")
// 	r.HandleFunc("/cache/clear", srv.ClearAllCaches).Methods("PUT")

// 	return r
// }

// func TestGetCache(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockCache := cache.NewMockCache(ctrl)
// 	mockCache.EXPECT().Get("testKey").Return("testValue", nil)

// 	router := setupRouter(mockCache, mockCache)

// 	req, err := http.NewRequest("GET", "/cache/testKey", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// }

// func TestSetCache(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockCache := cache.NewMockCache(ctrl)
// 	mockCache.EXPECT().Set("testKey", "testValue", gomock.Any()).Return(nil)

// 	router := setupRouter(mockCache, mockCache)

// 	var jsonStr = []byte(`{"key":"testKey","value":"testValue", "ttl":"15s"}`)
// 	req, err := http.NewRequest("POST", "/cache", bytes.NewBuffer(jsonStr))
// 	req.Header.Set("Content-Type", "application/json")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// }

// func TestDeleteCache(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockCache := cache.NewMockCache(ctrl)
// 	mockCache.EXPECT().Delete("testKey").Return(nil)

// 	router := setupRouter(mockCache, mockCache)

// 	req, err := http.NewRequest("DELETE", "/cache/testKey", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// }

// func TestGetWithTTL(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockCache := cache.NewMockCache(ctrl)
// 	mockCache.EXPECT().GetWithTTL("testKey").Return("testValue", time.Minute, nil)

// 	router := setupRouter(mockCache, mockCache)

// 	req, err := http.NewRequest("GET", "/cache/TTL/testKey", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// }
