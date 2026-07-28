package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func InjectCache(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cache", cache)
		c.Next()
	}
}
