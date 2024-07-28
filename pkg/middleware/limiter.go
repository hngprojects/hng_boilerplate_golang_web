package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(1, 4)
	return func(c *gin.Context) {

		if limiter.Allow() {
			c.Next()
		} else {
			rd := utility.BuildErrorResponse(http.StatusTooManyRequests, "error", "Limit exceed", nil, nil)
			c.JSON(http.StatusTooManyRequests, rd)
		}

	}
}
