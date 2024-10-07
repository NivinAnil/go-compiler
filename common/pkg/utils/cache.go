package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	client *redis.Client
	ctx    context.Context
}

// Create the interface for the cache client
type ICacheClient interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	Expire(key string, expiration time.Duration) error
	Flush() error
}

// NewCacheClient initializes a new CacheClient instance
func NewCacheClient(address string, password string, db int) *CacheClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &CacheClient{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Set a key with a value and an optional expiration time
func (cc *CacheClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := cc.client.Set(cc.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %v", key, err)
	}
	return nil
}

// Get a value for a given key
func (cc *CacheClient) Get(key string) (string, error) {
	fmt.Println("Attempting to get key:", key)

	val, err := cc.client.Get(cc.ctx, key).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s: %v", key, err)
	}

	fmt.Println("Value retrieved:", val)
	return val, nil
}

// Delete a key from the cache
func (cc *CacheClient) Delete(key string) error {
	err := cc.client.Del(cc.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %v", key, err)
	}
	return nil
}

// Check if a key exists in the cache
func (cc *CacheClient) Exists(key string) (bool, error) {
	exists, err := cc.client.Exists(cc.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %v", key, err)
	}
	return exists > 0, nil
}

// Set expiration for a key
func (cc *CacheClient) Expire(key string, expiration time.Duration) error {
	err := cc.client.Expire(cc.ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %v", key, err)
	}
	return nil
}

// Clear all keys in the cache (Flush)
func (cc *CacheClient) Flush() error {
	err := cc.client.FlushAll(cc.ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush cache: %v", err)
	}
	return nil
}
