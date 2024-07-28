package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPresent(t *testing.T) {
	redisUrl := "127.0.0.1:6379"

	cache := NewEventServerCache(redisUrl)

	cache.Add("token1")

	assert.True(t, cache.IsPresent("token1"))
	assert.True(t, !cache.IsPresent("token2"))

	cache.Add("token2")

	assert.True(t, cache.IsPresent("token2"))
}
