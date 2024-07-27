package cache

type CacheInterface interface {
	IsPresent(token string) bool
	Add(token string)
}