package auth

import (
	"context"
	"encoding/json"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthPayload struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

var auth_payload AuthPayload

// Claims struct for JWT
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var (
	Client_ID         string
	Client_Secret     string
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
	RedirectURL       = "http://127.0.0.1:8000/api/v1/auth/callback/google"
)

func init() {
	// Load the .env file
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Assigning Variables
	Client_ID = os.Getenv("CLIENT_ID")
	Client_Secret = os.Getenv("CLIENT_SECRET")

	// Initialize Google OAuth2 configuration
	googleOauthConfig = &oauth2.Config{
		//Redirect Must match Redirect URI in API Credentials
		RedirectURL:  RedirectURL,   // "http" used instead of "https" to resolve SSL certificate errors
		ClientID:     Client_ID,     // Your Google client ID
		ClientSecret: Client_Secret, // Your Google client secret
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func Handle_Google_Login(c *gin.Context) {
	// Generate the Google OAuth2 login URL with a state string for security
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	// Redirect the user to the Google login page
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func Handle_Google_Callback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid OAuth state"})
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "code exchange failed"})
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info"})
		return
	}

	// Logic for updating user in database
	user, err := updateUserInfo(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user info"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := generateJWT(user["email"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	// Respond with user info and tokens
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "User successfully authenticated",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})

}

// Update user information in database
func updateUserInfo(userInfo map[string]interface{}) (map[string]interface{}, error) {
	// Implement database update logic
	// This is a mock function; replace with actual database operations

	// Example: Check if user exists, if not, create new user
	// Update user details and sstore tokens

	// Mock user object
	user := map[string]interface{}{
		"id":          userInfo["id"],
		"email":       userInfo["email"],
		"name":        userInfo["name"],
		"given_name":  userInfo["given_name"],
		"family_name": userInfo["family_name"],
		"picture":     userInfo["picture"],
	}

	// Simulate saving to database
	return user, nil
}

func Handle_Token_Refresh(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err :=
		jwt.ParseWithClaims(
			request.RefreshToken,
			&Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	newAccessToken, _, err := generateJWT(claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

// ----------------------------Token Generation--------------------------------------------
// Define your secret key for signing the tokens
var jwtKey = []byte("your_secret_key")

// Generate JWT access and refresh tokens
func generateJWT(email string) (string, string, error) {
	// Define expiration times for the tokens
	accessTokenExpiration := time.Now().Add(30 * time.Minute)
	refreshTokenExpiration := time.Now().Add(7 * 24 * time.Hour)

	// Create the access token claims
	accessClaims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpiration.Unix(),
		},
	}
	// Create the JWT access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Create the refresh token claims
	refreshClaims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiration.Unix(),
		},
	}
	// Create the JWT refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
