package key

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

func CreateKey(db *gorm.DB, c *gin.Context) (gin.H, int, error) {
	userID, _ := middleware.GetIdFromToken(c)
	log.Print(userID)

	if userID == "" {
		return nil, http.StatusBadRequest, errors.New("User is not authenticated")
	}

	var existingKey models.Key
	if err := db.Where("user_id = ?", userID).First(&existingKey).Error; err == nil {
		return nil, http.StatusConflict, errors.New("Key for this user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "HNG_KIMIKO",
		AccountName: userID,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		db.Rollback()
		return nil, http.StatusNotFound, errors.New("User not found")
	}

	keyModel := models.Key{
		ID:     utility.GenerateUUID(),
		UserID: userID,
		Key:    secret.Secret(),
	}

	if err := db.Create(&keyModel).Error; err != nil {
		db.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	db.Commit()

	png, err := qrcode.Encode(secret.URL(), qrcode.Medium, 256)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return gin.H{
		"secret":  secret.Secret(),
		"qr_code": png,
	}, http.StatusCreated, nil
}

func VerifyKey(req models.VerifyKeyRequestModel, db *gorm.DB, c *gin.Context) (gin.H, int, error) {
	userID, _ := middleware.GetIdFromToken(c)
	key := req.Key
	if key == "" || userID == "" {
		return nil, http.StatusBadRequest, errors.New("Key and User ID are required")
	}
	var keyModel models.Key
	if err := db.Where("user_id = ?", userID).First(&keyModel).Error; err != nil {
		return nil, http.StatusNotFound, errors.New("Key not found")
	}

	if !totp.Validate(key, keyModel.Key) {
		return nil, http.StatusUnauthorized, errors.New("Invalid key")
	}

	return gin.H{
		"message": "Key verified successfully",
	}, http.StatusOK, nil
}
