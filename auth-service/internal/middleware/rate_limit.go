package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitMiddleware ограничивает кол-во запросов с одного IP
func RateLimitMiddleware(client *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// Ключ будем строить по IP клиента (или user_id, если авторизован)
		ip := c.ClientIP()
		key := "rate_limit:login:" + ip

		// Атомарно увеличиваем счетчик
		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit service error"})
			c.Abort()
			return
		}

		// Если это первый запрос, ставим время жизни ключа (TTL)
		if count == 1 {
			client.Expire(ctx, key, window)
		}

		// Если превысили лимит — блокируем
		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
