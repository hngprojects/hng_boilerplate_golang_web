package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Authorize(db *gorm.DB, inputRole ...models.RoleId) gin.HandlerFunc {
	// if no role is passed it would assume default user role
	return func(c *gin.Context) {

		var (
			tokenStr     string
			access_token models.AccessToken
		)

		bearerToken := c.GetHeader("Authorization")
		strArr := strings.Split(bearerToken, " ")
		if len(strArr) == 2 {
			tokenStr = strArr[1]
		}

		if tokenStr == "" {
			r := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Token could not be found!", "Unauthorized", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, r)
			return
		}

		token, err := TokenValid(tokenStr)
		if err != nil {
			r := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Token is invalid!", "Unauthorized", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, r)
			return
		}

		// access user claims

		claims := token.Claims.(jwt.MapClaims)

		// check if user id exists and fetch it
		userID, ok := claims["user_id"].(string) //convert the interface to string
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Token is invalid!", "Unauthorized", nil))
			return
		}

		// check if access id exists and fetch it
		accessID, ok := claims["access_uuid"].(string) //convert the interface to string
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Token is invalid!", "Unauthorized", nil))
			return
		}
		// check user session and also if token is valid in stored session

		access_token = models.AccessToken{ID: accessID}
		if code, err := access_token.GetByID(db); err != nil {
			c.AbortWithStatusJSON(code, utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Token is invalid!", "Unauthorized", nil))
			return
		}

		// check if session is valid

		if access_token.LoginAccessToken != tokenStr || userID != access_token.OwnerID || !access_token.IsLive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Session is invalid!", "Unauthorized", nil))
			return
		}

		// compare user role

		userRole := int(claims["role"].(float64)) //check if token is authorised for middleware
		var authorizedRole bool

		for _, role := range inputRole {
			if int(role) == userRole {
				authorizedRole = true
				break
			}
		}

		if !authorizedRole && len(inputRole) > 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.BuildErrorResponse(http.StatusUnauthorized, "error", "role not authorized!", "Unauthorized", nil))
			return
		}

		// check authorization status

		authoriseStatus, ok := claims["authorised"].(bool) //check if token is authorised for middleware
		if !ok && !authoriseStatus {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.BuildErrorResponse(http.StatusUnauthorized, "error", "status not authorized!", "Unauthorized", nil))
			return
		}

		// store user claims in Context
		// for accesiblity in controller

		c.Set("userClaims", claims)

		// call the next handler
		c.Next()

	}
}

func GetIdFromToken(c *gin.Context) (string, interface{}) {
	var tokenStr string
	bearerToken := c.GetHeader("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		tokenStr = strArr[1]
	}

	if tokenStr == "" {
		r := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Token could not be found!", "Unauthorized", nil)
		return "", r
	}

	token, err := TokenValid(tokenStr)
	if err != nil {
		r := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Token is invalid!", "Unauthorized", nil)
		return "", r
	}

	// access user claims

	claims := token.Claims.(jwt.MapClaims)
	id, ok := claims["user_id"].(string)
	if !ok {
		return "", utility.BuildErrorResponse(http.StatusForbidden, "error", "Forbidden", "Unauthorized", nil)
	}
	return id, ""
}