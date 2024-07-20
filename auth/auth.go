package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	Client_ID         string = "398923386395-0i6m9oll3046nc560dhcfbi5grj0blre.apps.googleusercontent.com"
	Client_Secret     string = "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"
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
)

func init() {
	// Initialize Google OAuth2 configuration
	googleOauthConfig = &oauth2.Config{
		//Redirect Must match Redirect URI in API Credentials
		RedirectURL:  "http://127.0.0.1:8000/api/v1/auth/callback/google", // "http" used instead of "https" to resolve SSL certificate errors
		ClientID:     Client_ID,                                           // Your Google client ID
		ClientSecret: Client_Secret,                                       // Your Google client secret
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
	// Verify the state string to protect against CSRF attacks
	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid OAuth state"})
		return
	}

	// Get the authorization code from the query parameters
	code := c.Query("code")
	// Exchange the authorization code for access and refresh tokens
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "code exchange failed"})
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	fmt.Println(response)
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

	//-----------------------------------
	//Logic for updating user in database
	//------------------------------------

	//Generating Tokens
	accessToken, refreshToken, err := generateJWT(userInfo["email"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	// Respond with the user info and tokens
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "User successfully authenticated",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userInfo,
	})

}

func Handle_Token_Refresh(c *gin.Context) {

}

// ----------------------------Token Generation--------------------------------------------
// Define your secret key for signing the tokens
var jwtKey = []byte("your_secret_key")

// Claims struct for JWT
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

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
