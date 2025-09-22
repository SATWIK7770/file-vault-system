package api

import (
	"backend/internal/service"
	"net/http"
	"os"
	"fmt"
	"sync"
    "time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read JWT from cookie instead of Authorization header
		tokenStr, err := c.Cookie("auth_token")
		if err != nil || tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth token"})
			c.Abort()
			return
		}

		// Validate token
		userID, err := service.ParseToken(tokenStr, []byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Store userID in context so handlers can use it
		c.Set("userID", uint(userID))

		fmt.Printf("DEBUG: File:%d", userID)
		c.Next()
	}
}


type RateLimiter struct {
    mu       sync.Mutex
    tokens   map[uint]int       // userID -> tokens left
    lastTime map[uint]time.Time // userID -> last refill
    limit    int                // calls per second
}

func NewRateLimiter(limit int) *RateLimiter {
    return &RateLimiter{
        tokens:   make(map[uint]int),
        lastTime: make(map[uint]time.Time),
        limit:    limit,
    }
}

func (rl *RateLimiter) RateMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetUint("userID")
        if userID == 0 {
            c.Next() // skip for unauthenticated
            return
        }

        rl.mu.Lock()
        defer rl.mu.Unlock()

        now := time.Now()
        last, exists := rl.lastTime[userID]
        if !exists {
            rl.tokens[userID] = rl.limit
            rl.lastTime[userID] = now
        }

        elapsed := now.Sub(last).Seconds()
        rl.tokens[userID] += int(elapsed * float64(rl.limit))
        if rl.tokens[userID] > rl.limit {
            rl.tokens[userID] = rl.limit
        }
        rl.lastTime[userID] = now

        if rl.tokens[userID] <= 0 {
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
            c.Abort()
            return
        }

        rl.tokens[userID]--
        c.Next()
    }
}
