package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iabdulzahid/golang_task_manager/internal/models"
)

var requestCount = make(map[string]int) // Tracks requests per IP
var mu sync.Mutex

// RateLimiter middleware that limits requests based on IP
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		ip := c.ClientIP()
		// Set a maximum of 100 requests per minute per IP
		maxRequests := 100
		window := 60 * time.Second

		// Initialize the request count if not already
		if requestCount[ip] == 0 {
			go resetRequestCount(ip, window)
		}

		// Check the number of requests from this IP
		if requestCount[ip] > maxRequests {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error: "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		// Increment request count
		requestCount[ip]++

		// Continue processing request
		c.Next()
	}
}

// Reset the request count after the window (60 seconds)
func resetRequestCount(ip string, window time.Duration) {
	time.Sleep(window)
	mu.Lock()
	defer mu.Unlock()
	requestCount[ip] = 0
}
