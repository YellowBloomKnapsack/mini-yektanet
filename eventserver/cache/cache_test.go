package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddPresent(t *testing.T) {
	redisUrl := "127.0.0.1:6379"
	redisExpireDuration := 500 * time.Millisecond

	cache := NewCacheService(redisExpireDuration, redisUrl)

	cache.Add("token1")

	assert.True(t, cache.IsPresent("token1"))
	assert.True(t, !cache.IsPresent("token2"))

	cache.Add("token2")

	assert.True(t, cache.IsPresent("token2"))
}

func TestExpiration(t *testing.T) {
	redisUrl := "127.0.0.1:6379"
	redisExpireDuration := 500 * time.Millisecond

	cache := NewCacheService(redisExpireDuration, redisUrl)

	cache.Add("token1")

	assert.True(t, cache.IsPresent("token1"))
	time.Sleep(redisExpireDuration + 200 * time.Millisecond)
	assert.True(t, !cache.IsPresent("token1"))
}
