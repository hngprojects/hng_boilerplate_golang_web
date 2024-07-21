package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func CreateToken(user models.User) (string, time.Time, error) {

	var (
		config  = config.GetConfig()
		UnixExp = time.Now().AddDate(0, 0, 7).Unix() // token valid for a week
		exp     = time.Now().AddDate(0, 0, 7)
	)

	//create token
	userid := user.ID
	userClaims := jwt.MapClaims{}

	// specify user claims
	userClaims["user_id"] = userid
	userClaims["exp"] = UnixExp
	userClaims["authorised"] = true

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	accessToken, err := token.SignedString([]byte(config.Server.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, exp, nil
}

// verify token

func verifyToken(tokenString string) (*jwt.Token, error) {
	config := config.GetConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Server.Secret), nil
	})
	if err != nil {
		return token, fmt.Errorf("Unauthorized")
	}
	return token, nil
}

// check if token is valid

func TokenValid(bearerToken string) (*jwt.Token, error) {
	token, err := verifyToken(bearerToken)
	if err != nil {
		if token != nil {
			return token, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Unauthorized")
	}
	return token, nil
}
