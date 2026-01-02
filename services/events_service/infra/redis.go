package infra

import (
	"os"
	"time"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func DefaultTTL() time.Duration {
	ttl, _ := time.ParseDuration(os.Getenv("REDIS_TTL_SECONDS") + "s")
	return ttl
}
