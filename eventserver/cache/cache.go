package cache

import (
	"context"
	"os"
	"log"
	"time"
	"fmt"

	"YellowBloomKnapsack/mini-yektanet/common/cache"
	"github.com/redis/go-redis/v9"
)

type EventServerCache struct {
	redisClient    *redis.Client
	resetPeriod 	time.Duration
	tableName 		string
}

func NewEventServerCache(redisUrl string, resetPeriod time.Duration) cache.CacheInterface {
	r := &EventServerCache{
		redisClient: redis.NewClient(&redis.Options{
			Addr:	  redisUrl,
        	Password: "", // no password set
        	DB:		  0,  // use default DB
		}),
		resetPeriod: resetPeriod,
	}

	r.reserveNewRedisTable()

	go r.bfResetService() // background service to reset bloom filter database 
	return r
}

func (r *EventServerCache) IsPresent(token string) bool {
	ctx := context.Background()
	// check with bloom filter
	exists, err := r.redisClient.Do(ctx, "BF.EXISTS", r.tableName, token).Bool()
	// exists, err := r.redisClient.BFExists(ctx, tableName(), token).Result()
	if err != nil {
		log.Println(err)
		return true
	}

	return exists
}

func (r *EventServerCache) Add(token string) {
	ctx := context.Background()

	// add to bloom filter
	_, err := r.redisClient.Do(ctx, "BF.ADD", r.bfTableName(), token).Bool()
	if err != nil {
		log.Println(err)
	}
}

func (r* EventServerCache) bfResetService() { // ask about older tables
	if r.resetPeriod != time.Duration(0) { // periodically reset bloom filter
		ticker := time.NewTicker(r.resetPeriod)
		for _ = range ticker.C {
			fmt.Println("reset redis.BF table on period ...")
			r.reserveNewRedisTable()
		}
	} else { // reset at the start of every day
		for {
			now := time.Now()
			nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			durationUntilReset := nextDay.Sub(now)
			time.Sleep(durationUntilReset)

			fmt.Println("reset redis.BF table at the start of the day ...")
			r.reserveNewRedisTable()
		}
	}
}

func (r *EventServerCache) reserveNewRedisTable() error {
	ctx := context.Background()

	initialSize := os.Getenv("REDIS_BF_INIT_SIZE") // initial capacity of the Bloom filter
	errorRate := os.Getenv("REDIS_BF_ERR_RATE")      // false positive rate

	// create new table name
	if r.resetPeriod != time.Duration(0) {
		r.tableName = time.Now().Format(time.RFC3339)
	} else {
		r.tableName = "bf_token_"+(time.Now().Format("2006.01.02"))
	}

	// store new table name in set of names for later removal
	r.redisClient.ZAdd(ctx, "bf_names_set", redis.Z{
		Score: 1,
		Member: r.tableName,
	})
	fmt.Println(r.tableName)

	_, err := r.redisClient.Do(ctx, "BF.RESERVE", r.tableName, errorRate, initialSize).Result()
	if err != nil && err.Error() != "ERR item exists" {
		log.Fatalf("Error creating Bloom filter: %v", err)
	}
	return err
}

func (r *EventServerCache) bfTableName() string {
	var tableName string
	if r.resetPeriod != time.Duration(0) {
		tableName = time.Now().Format(time.RFC3339)
	} else {
		tableName = "bf_token_"+(time.Now().Format("2006.01.02"))
	}
	fmt.Println(tableName)
	return tableName
}
