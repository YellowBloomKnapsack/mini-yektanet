package cache

import (
	"context"
	"log"
	"time"
	"strconv"
	"os"

	"YellowBloomKnapsack/mini-yektanet/common/cache"
	"github.com/redis/go-redis/v9"
)

type AdServerCache struct {
	redisClient *redis.Client
	brakeDuration time.Duration
}

func NewAdServerCache(redisUrl string) cache.CacheInterface {
	brakeSeconds, _ := strconv.Atoi(os.Getenv("BRAKE_DURATION_SECS"))
	return &AdServerCache{
		redisClient: redis.NewClient(&redis.Options{
			Addr:	  redisUrl,
        	Password: "", // no password set
        	DB:		  0,  // use default DB
		}),
		brakeDuration: time.Duration(brakeSeconds) * time.Second,
	}
}

func (r *AdServerCache) IsPresent(adID string) bool {
	ctx := context.Background()
	_, err := r.redisClient.Get(ctx, adID).Result()
	if err != nil {
		log.Println("could not search for " + adID + " on Adserver redis: " + err.Error())
	}
	// note: in the case of actual error in connection, we still return false - i.e. ad not braked
	return (err == redis.Nil)
}

func (r *AdServerCache) Add(adID string) {
	ctx := context.Background()
	err := r.redisClient.Set(ctx, adID, "", r.brakeDuration).Err()
	if err != nil {
		log.Println("could not insert " + adID + " in Adserver redis: " + err.Error())
	}
}