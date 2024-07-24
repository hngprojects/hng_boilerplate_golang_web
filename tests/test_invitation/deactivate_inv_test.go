package invitation_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/invite"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// GenerateTestJWT generates a JWT token for testing
func GenerateTestJWT() string {
	jwtKey := "test_secret"
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtKey))
	return tokenString
}

func DeactivateInvitation(ctx *gin.Context) {
	authHeader := ctx.GetHeader("authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
			"errors": []gin.H{
				{
					"field":   "authorization",
					"message": "User is not authorized to deactivate this invitation link",
				},
			},
			"status_code": 403,
		})
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
			"errors": []gin.H{
				{
					"field":   "authorization",
					"message": "Invalid authorization header",
				},
			},
			"status_code": 403,
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := middleware.TokenValid(tokenString)
	fmt.Print(token)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
			"errors": []gin.H{
				{
					"field":   "authorization",
					"message": err.Error(),
				},
			},
			"status_code": 403,
		})
		return
	}

	// Bind the request body to the invLink struct
	var invLink invite.InvLink
	if err := ctx.Bind(&invLink); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"errors": []gin.H{
				{
					"field":   "invitation_link",
					"message": "Invalid request body",
				},
			},
			"status_code": 400,
		})
		return
	}

	db := storage.Connection()
	var invitation models.Invitation
	result := db.Postgresql.Where("invitation_link = ?", invLink.InvitationLink).First(&invitation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Invitation link not found",
				"status_code": 404,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database error",
				"errors": []gin.H{
					{
						"field":   "database",
						"message": result.Error.Error(),
					},
				},
				"status_code": 500,
			})
		}
		return
	}

	// Update the isValid field to false
	invitation.IsValid = false
	if err := db.Postgresql.Save(&invitation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to deactivate invitation link",
			"errors": []gin.H{
				{
					"field":   "database",
					"message": err.Error(),
				},
			},
			"status_code": 500,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Invitation link deactivated successfully",
		"status_code": 200,
	})
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/deactivate-invitation", DeactivateInvitation) // Add route handler
	return r
}

func TestDeactivateInvitation(t *testing.T) {
	// Setup environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup in-memory SQLite database
	db := storage.Connection().Postgresql
	
	// Migrate the schema
	db.AutoMigrate(&models.Organisation{}, &models.Invitation{})
	storage.Connection().Postgresql = db

	// Create a test organization
	organization := models.Organisation{
		ID:   utility.GenerateUUID(),
		Name: "Test Organization",
	}
	db.Create(&organization)

	// Create a test invitation
	invitation := models.Invitation{
		ID:             utility.GenerateUUID(),
		Email:          "test@example.com",
		OrganizationID: organization.ID,
		IsValid:        true,
	}
	db.Create(&invitation)

	// Setup router
	router := SetupRouter()

	// Generate a valid JWT token
	token := GenerateTestJWT()

	// Test cases
	tests := []struct {
		name          string
		invitationLink string
		token         string
		expectedStatus int
		expectedValid  bool
	}{
		{
			name:           "Valid deactivation",
			invitationLink: invitation.ID,
			token:          token,
			expectedStatus: http.StatusOK,
			expectedValid:  false,
		},
		{
			name:           "Invalid token",
			invitationLink: invitation.ID,
			token:          "invalid_token",
			expectedStatus: http.StatusForbidden,
			expectedValid:  true,
		},
		{
			name:           "Invalid invitation link",
			invitationLink: "non_existing_link",
			token:          token,
			expectedStatus: http.StatusNotFound,
			expectedValid:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			body := invite.InvLink{
				InvitationLink: tc.invitationLink,
			}
			jsonBody, _ := json.Marshal(body)

			// Create HTTP request
			req, _ := http.NewRequest("POST", "/deactivate-invitation", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tc.token)

			// Create HTTP recorder
			w := httptest.NewRecorder()

			// Perform HTTP request
			router.ServeHTTP(w, req)

			// Assert HTTP status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Assert invitation validity
			var updatedInvitation models.Invitation
			db.First(&updatedInvitation, "id = ?", invitation.ID)
			assert.Equal(t, tc.expectedValid, updatedInvitation.IsValid)
		})
	}
}
