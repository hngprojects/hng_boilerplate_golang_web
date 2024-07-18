package router

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type DeactivateRequest struct {
	InvitationLink string `json:"invitation_link" binding:"required"`
}

func Setup(logger *utility.Logger, validator *validator.Validate, db *storage.Database, appConfiguration *config.App) *gin.Engine {
	if appConfiguration.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// Middlewares
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(config.GetConfig().Server.TrustedProxies)
	r.Use(middleware.Security())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.MaxMultipartMemory = 1 << 20 // 1MB

	// routers
	ApiVersion := "v2"
	Health(r, ApiVersion, validator, db, logger)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "HNGi Golang Boilerplate",
			"status":  http.StatusOK,
		})
	})

	r.POST("/api/v1/invite/deactivate", func(c *gin.Context) {
		var request DeactivateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Validation error",
				"errors":  []gin.H{{"field": "invitation_link", "message": "Invalid invitation link format"}},
				"status_code": 400,
			})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Forbidden",
				"errors":  []gin.H{{"field": "authorization", "message": "User is not authorized to deactivate this invitation link"}},
				"status_code": 403,
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Forbidden",
				"errors":  []gin.H{{"field": "authorization", "message": "Authorization header format must be Bearer {token}"}},
				"status_code": 403,
			})
			return
		}

		secretKey := []byte("your-secret-key")
		if err := utility.JWTValid(tokenString, secretKey); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Forbidden",
				"errors":  []gin.H{{"field": "authorization", "message": err.Error()}},
				"status_code": 403,
			})
			return
		}

		// Check if the invitation link exists and is valid in the database
		// This is a placeholder. Replace with your actual database check logic.
		invitationExists := true // Assume this checks the database

		if !invitationExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Validation error",
				"errors":  []gin.H{{"field": "invitation_link", "message": "Invalid or expired invitation link"}},
				"status_code": 400,
			})
			return
		}

		// Deactivate the invitation link
		// This is a placeholder. Replace with your actual database update logic.
		deactivated := true // Assume this updates the database

		if !deactivated {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
				"errors":  []gin.H{{"field": "deactivation", "message": "Failed to deactivate the invitation link"}},
				"status_code": 500,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Invitation link has been deactivated",
			"status_code": 200,
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		})
	})

	return r
}
