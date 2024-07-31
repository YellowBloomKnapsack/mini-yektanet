package kvstorage

type KVStorageInterface interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type KVStorage struct {
	storage map[string]string
}

func (kvs *KVStorage) Get(key string) (string, error) {
	return kvs.storage[key], nil
}

func (kvs *KVStorage) Set(key string, value string) error {
	kvs.storage[key] = value
	return nil
}
