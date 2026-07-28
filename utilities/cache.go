package utilities

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func NewCache(defaultExpiration, cleanupInterval time.Duration) *cache.Cache {
	return cache.New(defaultExpiration, cleanupInterval)
}
