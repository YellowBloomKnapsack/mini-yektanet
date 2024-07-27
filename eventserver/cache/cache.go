package cache

import (
	"context"
	"os"
	"log"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/cache"
	"github.com/redis/go-redis/v9"
)

type EventServerCache struct {
	redisClient    *redis.Client
}

func NewEventServerCache(redisUrl string) cache.CacheInterface {
	r := &EventServerCache{
		redisClient: redis.NewClient(&redis.Options{
			Addr:	  redisUrl,
        	Password: "", // no password set
        	DB:		  0,  // use default DB
		}),
	}

	ctx := context.Background()
	initialSize := os.Getenv("REDIS_BF_INIT_SIZE") // initial capacity of the Bloom filter
	errorRate := os.Getenv("REDIS_BF_ERR_RATE") // false positive rate
	_, err := r.redisClient.Do(ctx, "BF.RESERVE", bfTableName(), errorRate, initialSize).Result()
	if err != nil {
		log.Printf("Error recreating Bloom filter: %v", err)
	}
	go bfResetService(r) // background service to reset bloom filter database 
	return r
}

func (r *EventServerCache) IsPresent(token string) bool {
	ctx := context.Background()
	// check with bloom filter
	// exists, err := r.redisClient.Do(ctx, "BF.EXISTS", bfTableName(), token).Bool()
	exists, err := r.redisClient.BFExists(ctx, bfTableName(), token).Result()
	if err != nil {
		log.Println(err)
		return true
	}

	return exists
}

func (r *EventServerCache) Add(token string) {
	ctx := context.Background()

	// add to bloom filter
	_, err := r.redisClient.Do(ctx, "BF.ADD", bfTableName(), token).Bool()
	if err != nil {
		log.Println(err)
	}
}

func bfResetService(r *EventServerCache) {
	for {
		now := time.Now()
		nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilReset := nextDay.Sub(now)
		time.Sleep(durationUntilReset)

		ctx := context.Background()
		
		initialSize := os.Getenv("REDIS_BF") // initial capacity of the Bloom filter
		errorRate := 0.01      // false positive rate

		_, err := r.redisClient.Do(ctx, "BF.RESERVE", bfTableName(), errorRate, initialSize).Result()
		if err != nil {
			log.Fatalf("Error recreating Bloom filter: %v", err)
		}
	}
}