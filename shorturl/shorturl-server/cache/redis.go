package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)
import pkg_redis "shorturl/pkg/db/redis"

type RedisKVCache struct {
	redisClient *redis.Client
	destroy     func()
}

func newRedisKVCache(redisClient *redis.Client, destroy func()) *RedisKVCache {
	return &RedisKVCache{
		redisClient: redisClient,
		destroy:     destroy,
	}
}
func getKey(key string) string {
	return pkg_redis.GetKey(key)
}
func (c *RedisKVCache) Get(key string) (string, error) {
	key = getKey(key)
	res, err := c.redisClient.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return res, err
}

func (c *RedisKVCache) Set(key string, value string, ttl int) error {
	key = getKey(key)
	err := c.redisClient.Set(context.Background(), key, value, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		log.Printf("Error setting value in Redis: %v", err)
		return err
	}
	return nil
}

func (c *RedisKVCache) Destroy() {
	if c.destroy != nil {
		c.destroy()
	}
}

// redis factory
type redisCacheFactory struct {
	redisPool pkg_redis.RedisPool
}

func NewRedisCacheFactory(redisPool pkg_redis.RedisPool) CacheFactory {
	return &redisCacheFactory{
		redisPool: redisPool,
	}
}
func (f *redisCacheFactory) NewKVCache() KVCache {
	client := f.redisPool.Get()
	return newRedisKVCache(client, func() {
		f.redisPool.Put(client)
	})
}
