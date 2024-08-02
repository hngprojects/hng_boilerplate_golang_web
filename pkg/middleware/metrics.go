package middleware

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
)

var (
	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status", "app"})
)

func init() {
	// Register the custom metric with the default registry
	prometheus.MustRegister(httpRequestDuration)
}

// Metrics is a middleware function for tracking HTTP request durations
func Metrics(cfg *config.Configuration) gin.HandlerFunc {
	appName := cfg.App.Name
	if appName == "" {
		log.Println("Warning: APP_NAME is not set in configuration")
		appName = "unknown_app"
	}

	return func(c *gin.Context) {
		// Start a timer to track request duration
		startTime := time.Now()

		// Get the request method
		method := c.Request.Method

		// Allow the request to be processed by calling the next handler
		c.Next()

		// Get the response status code after the request is processed
		statusCode := c.Writer.Status()

		// Calculate and record the request duration
		elapsed := time.Since(startTime).Seconds()
		httpRequestDuration.WithLabelValues(method, c.FullPath(), strconv.Itoa(statusCode), appName).Observe(elapsed)
	}
}
