package cache

import (
	"os"
	"raise-child/constants/env"

	"github.com/redis/go-redis/v9"
)

var _client *redis.Client

func InitializeRedisClient() *redis.Client {
	if _client == nil {
		_client = redis.NewClient(&redis.Options{
			Addr:     os.Getenv(env.REDIS_ADDRESS),
			Password: os.Getenv(env.REDIS_PASSWORD),
			DB:       0,
		})
	}

	return _client
}
