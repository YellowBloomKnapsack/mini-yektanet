package kvstorage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type KVStorageInterface interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type KVStorage struct {
	client *redis.Client
}

// NewKVStorage creates a new KVStorage instance with a connected Redis client
func NewKVStorage(addr string) *KVStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	return &KVStorage{client: client}
}

func (kvs *KVStorage) Get(key string) (string, error) {
	val, err := kvs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("could not find key %s", key) // Key does not exist
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (kvs *KVStorage) Set(key string, value string) error {
	err := kvs.client.Set(ctx, key, value, 0).Err()
	return err
}
