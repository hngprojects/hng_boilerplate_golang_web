package middleware

import (
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// GzipWithExclusion applies gzip compression conditionally based on the request path.
func GzipWithExclusion(excludedPaths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range excludedPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}
		gzip.Gzip(gzip.DefaultCompression)(c)
	}
}
