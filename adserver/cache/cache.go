package cache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/cache"
	"github.com/redis/go-redis/v9"
)

type AdServerCache struct {
	redisClient   *redis.Client
	brakeDuration time.Duration
	prefixStore   string
}

func NewAdServerCache(redisUrl string) cache.CacheInterface {
	brakeSeconds, _ := strconv.Atoi(os.Getenv("BRAKE_DURATION_SECS"))
	return &AdServerCache{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisUrl,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		brakeDuration: time.Duration(brakeSeconds) * time.Second,
		prefixStore:   os.Getenv("BRAKE_IDS_REDIS_PREFIX"),
	}
}

func (r *AdServerCache) IsPresent(key string) bool {
	ctx := context.Background()
	present, err := r.redisClient.Exists(ctx, r.prefixStore+":"+key).Result()
	if err != nil {
		log.Print(err)
		return false
	}

	return bool(present == 1)
}

func (r *AdServerCache) Add(key string) {
	ctx := context.Background()
	err := r.redisClient.Set(ctx, r.prefixStore+":"+key, "", r.brakeDuration).Err()
	if err != nil {
		log.Println("could not insert " + key + " in adserver redis: " + err.Error())
	}
}
