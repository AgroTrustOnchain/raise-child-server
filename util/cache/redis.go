package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/shared"
	"raise-child/util"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// if _client == nil {
// 	_client = redis.NewClient(&redis.Options{
// 		Addr:         os.Getenv(env.REDIS_ADDRESS),
// 		Password:     os.Getenv(env.REDIS_PASSWORD),
// 		PoolSize:     10,
// 		MinIdleConns: 3,
// 		DB:           0,
// 	})
// }

// if _redisCache == nil {
// 	_redisCache = &redisCache{
// 		client: _client,
// 	}
// }

type redisCache struct {
	client    *redis.Client
	errLogger *log.Logger
}

type IRedisCache interface {
	Set(key string, value any, duration time.Duration, ctx context.Context)
	Get(key string, value any, ctx context.Context) bool
	Delete(key string, ctx context.Context)
}

var (
	_redisCache *redisCache
	_once       sync.Once
)

func InitializeRedisCache() IRedisCache {
	_once.Do(func() {
		options, _ := redis.ParseURL(os.Getenv(env.REDIS_URL))

		options.PoolSize = 10
		options.MinIdleConns = 2

		options.DialTimeout = 5 * time.Second
		options.ReadTimeout = 3 * time.Second
		options.WriteTimeout = 3 * time.Second

		options.MaxRetries = 3
		options.MinRetryBackoff = 500 * time.Millisecond
		options.MaxRetryBackoff = 2 * time.Second

		_redisCache = &redisCache{
			client:    redis.NewClient(options),
			errLogger: util.GetLogConfig(shared.ERROR_LEVEL),
		}
	})

	return _redisCache
}

// Delete implements IRedisCache.
func (r *redisCache) Delete(key string, ctx context.Context) {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		r.errLogger.Println(err)
	}
}

// Get implements IRedisCache.
func (r *redisCache) Get(key string, value any, ctx context.Context) bool {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.errLogger.Println(err)
		}

		return false
	}

	return json.Unmarshal(data, value) == nil
}

// Set implements IRedisCache.
func (r *redisCache) Set(key string, value any, duration time.Duration, ctx context.Context) {
	data, err := json.Marshal(value)
	if err != nil {
		r.errLogger.Println(err)
		return
	}

	if err := r.client.Set(ctx, key, data, duration).Err(); err != nil {
		r.errLogger.Println(err)
	}
}
