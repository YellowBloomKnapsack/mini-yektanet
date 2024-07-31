package main

import (
	"YellowBloomKnapsack/mini-yektanet/adserver/kvstorage"
	_ "fmt"
	"time"
)

func main() {
	var kvstorage kvstorage.KVStorageInterface

	crawler := NewCrawler(kvstorage)

	ticker := time.NewTicker(1 * time.Hour)

	crawler.Crawl()

	go func() {
		for {
			select {
			case <-ticker.C:
				crawler.Crawl()
			}
		}
	}()

	select {}
}
