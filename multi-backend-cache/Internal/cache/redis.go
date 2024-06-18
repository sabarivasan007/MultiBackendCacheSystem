// Internal/cache/redis_cache.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	utils "multi-backend-cache/packageUtils/Utils"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// var NotFound = errors.New("key does not exist")

func NewRedisCache(addr string, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	// logrus.Info("Redis initialized")
	// // Print the defaults
	// logrus.Infof("Default PoolSize: %d", client.Options().PoolSize)
	// logrus.Infof("Default MinIdleConns: %d", client.Options().MinIdleConns)
	// logrus.Infof("Default MaxRetries: %d", client.Options().MaxRetries)
	// logrus.Infof("Default DialTimeout: %s", client.Options().DialTimeout)
	// logrus.Infof("Default ReadTimeout: %s", client.Options().ReadTimeout)
	// logrus.Infof("Default WriteTimeout: %s", client.Options().WriteTimeout)
	return &RedisCache{client: client, ttl: ttl}
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			// return nil, fmt.Errorf("key does not exist")
			return nil, utils.NotFound
		}
		return nil, err
	}
	var data interface{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// func (r *RedisCache) GetWithTTL(key string) (interface{}, time.Duration, time.Time, error) {
// 	ctx := context.Background()

// 	value, err := r.client.Get(ctx, key).Result()
// 	if err != nil {
// 		return nil, 0, time.Time{}, err
// 	}

// 	ttl, err := r.client.TTL(ctx, key).Result()
// 	if err != nil {
// 		return nil, 0, time.Time{}, err
// 	}
// 	fmt.Printf("%T\n", ttl)
// 	fmt.Println("------------------------",ttl)
// 	expiryTime := CalculateExpiryTime(ttl)

// 	var data interface{}
// 	err = json.Unmarshal([]byte(value), &data)
// 	if err != nil {
// 		return nil, 0, time.Time{}, fmt.Errorf("failed to unmarshal value: %v", err)
// 	}

// 	return data, ttl, expiryTime, nil
// }

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {

	ttlDuration := time.Duration(ttl) * time.Second
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	//Use the default TTL if the provided ttl is not provided
	actualTTL := ttlDuration
	if ttl <= 0 {
		actualTTL = r.ttl
	}
	// fmt.Printf("%T\n", actualTTL)
	// fmt.Println("----------------", actualTTL)

	fmt.Printf("Setting KEY: %s with VALUE: %s and TTL: %v seconds\n", key, value, actualTTL)
	return r.client.Set(context.Background(), key, val, actualTTL).Err()
}

func (r *RedisCache) Delete(key string) error {
	result, err := r.client.Del(context.Background(), key).Result()
	if err != nil {
		logrus.Errorf("Delete: error deleting key %s: %v", key, err)
		return err
	}

	if result == 0 {
		// logrus.Errorf("Delete: key %s not found", key)
		return utils.NotFound
	}
	return nil
}

func (r *RedisCache) Clear() error {
	return r.client.FlushDB(context.Background()).Err()
}
