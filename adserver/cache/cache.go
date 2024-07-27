package cache

import (
	"context"
	"fmt"
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
	}
}

func (r *AdServerCache) IsPresent(token string) bool {
	fmt.Println(token)
	ctx := context.Background()
	present, err := r.redisClient.Exists(ctx, token).Result()
	if err != nil {
		log.Print(err)
		return false
	}

	return bool(present == 1)
}

func (r *AdServerCache) Add(adID string) {
	ctx := context.Background()
	err := r.redisClient.Set(ctx, adID, "", r.brakeDuration).Err()
	if err != nil {
		log.Println("could not insert " + adID + " in Adserver redis: " + err.Error())
	}
}
