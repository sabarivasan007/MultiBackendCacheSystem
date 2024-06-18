package cache

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	utils "multi-backend-cache/packageUtils/Utils"
	"reflect"
	"sync"
	"time"

	"multi-backend-cache/Internal/config"

	"github.com/pbnjay/memory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// type MyDuration = time.Duration

type CacheData struct { //node

	Key        string        `json:"key" example:"1"`
	Value      interface{}   `json:"value" `
	TTL        time.Duration `json:"ttl" example:"100"`
	ExpiryTime time.Time     `json:"expirytime" example:"2021-05-25T00:53:16.535668Z" format:"date-time" swaggerignore:"true"`
}

// LRUCache represents the LRU cache, that consists of capacity, linkedlist as list, hashmap as index and lock
type LRUCache struct {
	capacity   int // in bytes
	used       int // in bytes
	list       *list.List
	index      map[string]*list.Element //key-> sring, value -> pointer to the list(*list.Element)
	lock       sync.Mutex
	defaultTTL time.Duration
}

type FixedTenantsCaches struct {
	caches map[string]*LRUCache
}

// GetCache retrieves the cache for the specified tenant.
func (ftc *FixedTenantsCaches) GetCache(tenantID string) *LRUCache {
	if cache, exists := ftc.caches[tenantID]; exists {
		return cache
	}
	return nil // Optionally handle the case where tenantID is not recognized.
}

// Initializes fixed tenant caches with predefined capacities.
func NewFixedTenantsCaches(totalCacheMemory int, defaultTTL time.Duration) *FixedTenantsCaches {
	// totalCacheMemory := float64(memory.TotalMemory()) * 0.15
	// defaultTTL := viper.GetInt("defaultTTL")
	tenantCaches := make(map[string]*LRUCache)
	tenantIDs := config.AppConfig.TenantIDs
	fmt.Println("tenantIDs", tenantIDs)
	eachTenantMemory := int(totalCacheMemory) / len(tenantIDs)
	//ttl := time.Duration(defaultTTL) * time.Second
	for _, id := range tenantIDs {
		tenantCaches[id] = NewLRUCache(eachTenantMemory, defaultTTL) // Define capacity per tenant here.
	}
	return &FixedTenantsCaches{
		caches: tenantCaches,
	}
}

// Checks the cache is expired or not
func IsExpired(expiryTime time.Time) bool {
	return time.Now().After(expiryTime)
}

// NewLRUCache creates a new LRU cache with the given capacity and ttl
func NewLRUCache(capacity int, defaultTTL time.Duration) *LRUCache {
	lru := &LRUCache{ // The "&" operator returns a pointer to the newly created LRUCache instance.
		capacity:   capacity,
		list:       list.New(),
		index:      make(map[string]*list.Element),
		defaultTTL: defaultTTL,
	}
	go DeleteExpiredCache(lru)

	// Go-routine that runs every one hour and checks the system memory and calculate cache memory
	// if len(tenantCaches) != 0 {
	// 	go checkMemoryForTenants(tenantCaches)
	// } else {
	go checkMemory(lru)

	// }
	return lru
}

// Go-routine that runs concurrently and for each 5 seconds scans the memory and deletes the
// expired ones
func DeleteExpiredCache(lru *LRUCache) {
	for range time.Tick(5 * time.Second) {
		lru.lock.Lock()
		// Iterate over the cache items and delete expired ones.
		for key, item := range lru.index {
			if IsExpired(item.Value.(*CacheData).ExpiryTime) {
				delete(lru.index, key)
				lru.list.Remove(item)
				log.Println("Deleted cache key ", key, item.Value.(*CacheData).ExpiryTime)
				size := CalculateSize(item.Value.(*CacheData))
				lru.used -= size
			}
		}
		lru.lock.Unlock()
	}
}

func checkMemory(lru *LRUCache) {
	for range time.Tick(1 * time.Hour) {
		lru.lock.Lock()
		cacheMemory := float64(memory.TotalMemory()) * viper.GetFloat64("MemoryUsagePercentage")
		fmt.Println("cache", cacheMemory)
		lru.capacity = int(cacheMemory)
		lru.lock.Unlock()
	}
}

// func checkMemoryForTenants(tenantCaches map[string]*LRUCache) {
// 	numberOfTenants := len(tenantCaches)
// 	for range time.Tick(1 * time.Hour) {
// 		cacheMemory := float64(memory.TotalMemory()) * 0.15
// 		fmt.Println("cahememory", int(cacheMemory))
// 		for tenant, cache := range tenantCaches {
// 			cache.lock.Lock()
// 			cache.capacity = int(cacheMemory) / numberOfTenants
// 			log.Println("Latest cache memory :: ", cache.capacity, tenant)
// 			cache.lock.Unlock()
// 		}

// 	}
// }

func CalculateExpiryTime(ttl time.Duration) time.Time {
	logrus.Debug("Calculating Expiry Time......")
	//return time.Now().Add(ttl)
	return time.Now().Add(ttl * time.Second)

}

// GetAllCache retrieves all values from the cache
func (c *LRUCache) GetAllCache() []*CacheData {
	c.lock.Lock()
	defer c.lock.Unlock()

	var allCacheData []*CacheData
	for element := c.list.Front(); element != nil; element = element.Next() {
		node := element.Value.(*CacheData)
		nodeJSON, err := json.Marshal(node)
		if err != nil {
			logrus.Error("Error marshalling node to JSON:", err)
		} else {
			logrus.Debugf("All cached data without Expired Cache: %s", string(nodeJSON))
		}
		if IsExpired(node.ExpiryTime) {
			// Entry has expired, remove it
			c.list.Remove(element)
			delete(c.index, node.Key)
			c.used -= CalculateSize(node)
		} else {
			allCacheData = append(allCacheData, node)
		}
	}
	return allCacheData
}

// GetCache retrieves a value from the cache
func (c *LRUCache) Get(key string) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var Err error
	if element, found := c.index[key]; found {
		logrus.Debugf("Existing cache found for key %s: %v", key, element.Value)
		node := element.Value.(*CacheData)
		nodeJSON, err := json.Marshal(node)
		if err != nil {
			logrus.Error("Error marshalling node to JSON:", err)
		} else {
			logrus.Infof("Cache data for key %s: %s", key, string(nodeJSON))
		}
		size := CalculateSize(node)
		if IsExpired(node.ExpiryTime) { // Check if the entry has expired
			c.list.Remove(element)
			delete(c.index, key)
			c.used -= size
			// return nil, fmt.Errorf("key %s has expired", key)
			return nil, utils.NotFound
		}
		log.Println("stp-1 node.TTL:", node.TTL.Seconds())
		// expiryTime := c.CalculateExpiryTime(node.TTL) //extend expiry
		log.Println("stp-2 node.TTL:", node.TTL)
		//node.TTL = time.Duration((node.TTL * time.Nanosecond).Seconds())
		log.Println("stp-3 after node.TTL:", node.TTL)
		// node.ExpiryTime = expiryTime
		c.list.MoveToFront(element)
		return node.Value, Err
	} else {
		return nil, utils.NotFound
	}
}

// Get cache along with data
// func (c *LRUCache) GetWithTTL(key string) (interface{}, time.Duration, time.Time, error) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	var Err error
// 	if element, found := c.index[key]; found {
// 		log.Println("Exsiting cache found", element.Value)
// 		node := element.Value.(*CacheData)
// 		nodeJSON, err := json.Marshal(node)
// 		if err != nil {
// 			log.Println("Error marshalling node to JSON:", err)
// 		} else {
// 			log.Printf("Cache data for key %s is: %s", key, string(nodeJSON))
// 		}
// 		size := CalculateSize(node)
// 		if IsExpired(node.ExpiryTime) { // Check if the entry has expired
// 			c.list.Remove(element)
// 			delete(c.index, key)
// 			c.used -= size
// 			// return nil, err
// 		}
// 		log.Println("stp-1 node.TTL:", node.TTL.Seconds())
// 		// expiryTime := CalculateExpiryTime(node.TTL) //extend expiry
// 		log.Println("stp-2 node.TTL:", node.TTL)
// 		//node.TTL = time.Duration((node.TTL * time.Nanosecond).Seconds())
// 		log.Println("stp-3 after node.TTL:", node.TTL)
// 		// node.ExpiryTime = expiryTime
// 		c.list.MoveToFront(element)
// 		ttl := time.Until(node.ExpiryTime)
// 		return node.Value, ttl, node.ExpiryTime, Err
// 	} else {
// 		Err := errors.New("key not found")
// 		return Err, 0, time.Time{}, nil
// 	}
// }

// setCache adds a value to the cache or updates the exisiting value
func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) error {
	logrus.Debugf("Setting key %s", key)
	c.lock.Lock()
	defer c.lock.Unlock()
	if ttl <= 0 {
		ttl = c.defaultTTL
	}
	logrus.Debugf("TTL for key %s: %s", key, ttl)
	// if key == "" {
	// 	err := errors.New("key must not be null")
	// 	logrus.Error("Set error:", err)
	// 	return err
	// }
	expiryTime := CalculateExpiryTime(ttl)
	if element, found := c.index[key]; found {
		logrus.Infof("Updating existing cache for key %s", key)
		c.list.MoveToFront(element)
		node := element.Value.(*CacheData)
		oldSize := CalculateSize(node)
		c.used -= oldSize
		fmt.Println("ttl in update", ttl, expiryTime)
		node.Value = value
		node.TTL = ttl
		node.ExpiryTime = expiryTime
		newSize := CalculateSize(node)
		c.used += newSize
		log.Println("TTL in node :: ", node.TTL)
		// return node
	} else {
		logrus.Infof("Creating new cache node for key %s", key)
		fmt.Println("ttl in update", ttl, expiryTime)
		newNode := &CacheData{Key: key, Value: value, TTL: ttl, ExpiryTime: expiryTime}
		fmt.Println("newNode", newNode)
		fmt.Println("ttl while setting cache :: ", newNode.TTL)
		size := CalculateSize(newNode)

		for c.used+size > c.capacity { // Recursively checked for freeing the last element to store the new one.
			logrus.Warn("Capacity Exceeded. Removing least recently used items.")
			backElement := c.list.Back()
			if backElement != nil {
				backNode := backElement.Value.(*CacheData)
				backnodeSize := CalculateSize(backNode)
				delete(c.index, backNode.Key) // map delete
				c.list.Remove(backElement)    // node delete
				c.used -= backnodeSize
			}
		}
		element := c.list.PushFront(newNode)
		c.index[key] = element
		c.used += size
	}
	return nil
}

// DeleteCache deletes a value from the cache
func (c *LRUCache) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, found := c.index[key]; found {
		node := element.Value.(*CacheData)
		size := CalculateSize(node)
		c.list.Remove(element)
		delete(c.index, node.Key)
		c.used -= size
		logrus.Infof("Deleted cache for key %s", key)
		return nil
	} else {
		return utils.NotFound
	}
}

func CalculateSize(node *CacheData) int {
	size := reflect.TypeOf(*node).Size()
	return int(size)
}

func (c *LRUCache) Clear() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.list.Init()
	c.index = make(map[string]*list.Element)
	c.used = 0
	return nil
}
