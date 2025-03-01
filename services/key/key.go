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

// KeyService defines the interface for key-related operations
type KeyService interface {
	CreateKey(c *gin.Context) (gin.H, int, error)
	VerifyKey(req models.VerifyKeyRequestModel, c *gin.Context) (gin.H, int, error)
}

// keyServiceImpl is the concrete implementation of KeyService
type keyServiceImpl struct {
	db *gorm.DB
}

// NewKeyService initializes a new KeyService instance
func NewKeyService(db *gorm.DB) KeyService {
	return &keyServiceImpl{db: db}
}

// CreateKey generates a new OTP key for a user
func (s *keyServiceImpl) CreateKey(c *gin.Context) (gin.H, int, error) {
	userID, _ := middleware.GetIdFromToken(c)
	log.Print(userID)

	if userID == "" {
		return nil, http.StatusBadRequest, errors.New("user is not authenticated")
	}

	var existingKey models.Key
	if err := s.db.Where("user_id = ?", userID).First(&existingKey).Error; err == nil {
		return nil, http.StatusConflict, errors.New("key for this user already exists")
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
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	keyModel := models.Key{
		ID:     utility.GenerateUUID(),
		UserID: userID,
		Key:    secret.Secret(),
	}

	if err := s.db.Create(&keyModel).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	png, err := qrcode.Encode(secret.URL(), qrcode.Medium, 256)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return gin.H{
		"secret":  secret.Secret(),
		"qr_code": png,
	}, http.StatusCreated, nil
}

// VerifyKey checks if a given OTP key is valid
func (s *keyServiceImpl) VerifyKey(req models.VerifyKeyRequestModel, c *gin.Context) (gin.H, int, error) {
	userID, _ := middleware.GetIdFromToken(c)
	key := req.Key
	if key == "" || userID == "" {
		return nil, http.StatusBadRequest, errors.New("key and user ID are required")
	}

	var keyModel models.Key
	if err := s.db.Where("user_id = ?", userID).First(&keyModel).Error; err != nil {
		return nil, http.StatusNotFound, errors.New("key not found")
	}

	if !totp.Validate(key, keyModel.Key) {
		return nil, http.StatusUnauthorized, errors.New("invalid key")
	}

	return gin.H{
		"message": "key verified successfully",
	}, http.StatusOK, nil
}
