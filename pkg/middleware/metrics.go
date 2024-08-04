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
	}, []string{"method", "route", "status", "name", "interpreter"})
)

func init() {
	prometheus.MustRegister(httpRequestDuration)
}

func Metrics(cfg *config.Configuration) gin.HandlerFunc {
	appName := cfg.App.Name
	if appName == "" {
		log.Println("Warning: APP_NAME is not set in configuration")
		appName = "unknown_app"
	}

	return func(c *gin.Context) {
		startTime := time.Now()

		method := c.Request.Method

		c.Next()

		statusCode := c.Writer.Status()

		elapsed := time.Since(startTime).Seconds()
		httpRequestDuration.WithLabelValues(method, c.FullPath(), strconv.Itoa(statusCode), appName, "none").Observe(elapsed)
	}
}
