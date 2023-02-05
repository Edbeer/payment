package red

import (
	"time"

	"github.com/go-redis/redis/v9"
)

func NewRedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:         ":6379",
		MinIdleConns: 200,
		PoolSize:     12000,
		PoolTimeout:  time.Duration(240) * time.Second,
		// Password:     "", // no password set
		DB:           0, // use default DB
	})

	return client
}