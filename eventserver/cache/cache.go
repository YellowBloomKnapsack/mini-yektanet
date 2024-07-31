package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/cache"
	"YellowBloomKnapsack/mini-yektanet/eventserver/grafana"

	"github.com/redis/go-redis/v9"
)

type EventServerCache struct {
	redisClient    *redis.Client
	tableName 		string
}

func NewEventServerCache(redisUrl string) cache.CacheInterface {
	r := &EventServerCache{
		redisClient: redis.NewClient(&redis.Options{
			Addr:	  redisUrl,
        	Password: "", // no password set
        	DB:		  0,  // use default DB
		}),
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
		grafana.CacheOperationsTotal.WithLabelValues("check", "error").Inc()
		return true
	}

	if exists {
        grafana.CacheOperationsTotal.WithLabelValues("check", "hit").Inc()
    } else {
        grafana.CacheOperationsTotal.WithLabelValues("check", "miss").Inc()
    }


	return exists
}

func (r *EventServerCache) Add(token string) {
	ctx := context.Background()

	// add to bloom filter
	_, err := r.redisClient.Do(ctx, "BF.ADD", r.tableName, token).Bool()
	if err != nil {
        log.Println(err)
        grafana.CacheOperationsTotal.WithLabelValues("add", "error").Inc()
    } else {
        grafana.CacheOperationsTotal.WithLabelValues("add", "success").Inc()
    }

}

func (r* EventServerCache) bfResetService() {
	for { // reset at the start of every day
		now := time.Now()
		nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilReset := nextDay.Sub(now)
		time.Sleep(durationUntilReset)

		fmt.Println("reset redis.BF table at the start of the day ...")
		r.reserveNewRedisTable()
	}
}

func (r *EventServerCache) reserveNewRedisTable() error {
	ctx := context.Background()

	initialSize := os.Getenv("REDIS_BF_INIT_SIZE") // initial capacity of the Bloom filter
	errorRate := os.Getenv("REDIS_BF_ERR_RATE")      // false positive rate

	// create new table name
	r.tableName = "bf_token_"+(time.Now().Format("2006.01.02"))

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
