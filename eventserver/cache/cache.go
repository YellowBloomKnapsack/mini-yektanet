package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheInterface interface {
	IsPresent(token string) bool
	Add(token string)
}

type CacheService struct {
	redisClient    redis.Client
	expireDuration time.Duration
}

func NewCacheService(expireDuration time.Duration, redisUrl string) CacheInterface {
	return &CacheService{
		expireDuration: expireDuration,
		redisClient: *redis.NewClient(&redis.Options{
			Addr:     redisUrl,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (r *CacheService) IsPresent(token string) bool {
	ctx := context.Background()
	present, err := r.redisClient.Exists(ctx, token).Result()
	if err != nil {
		log.Print(err)
		return false
	}

	return (present == 1)
}

func (r *CacheService) Add(token string) {
	ctx := context.Background()
	err := r.redisClient.Set(ctx, token, true, r.expireDuration).Err()

	if err != nil {
		log.Print(err)
	}
}
