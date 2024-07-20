package invitations

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"gorm.io/gorm"
)

type Invitation struct {
	ID             uint   `gorm:"primaryKey"`
	InvitationLink string `json:"invitation_link"`
	IsValid        bool   `json:"is_valid"`
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	jwtKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func CreateInvitation(ctx *gin.Context) {
	// Implementation for creating an invitation
}

type InvLink struct {
	InvitationLink string `json:"invitation_link"`
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
	token, err := VerifyJWT(tokenString)
	fmt.Println(token)
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
	var invLink InvLink
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

	// Check if the invitation link exists in the database
	db := storage.Connection()
	var invitation Invitation 
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
		"message": "Invitation link deactivated successfully",
		"status_code": 200,
	})
}
