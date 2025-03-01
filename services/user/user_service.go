package user

import (
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(userIDStr string) (models.User, int, error)
	GetUserByEmail(email string) (models.User, error)
	GetAUser(userIDStr string, requesterID string) (*models.User, int, error)
	GetAUserOrganisation(targetUserID string, requesterUserID string) (*[]models.Organisation, int, error)
	DeleteAUser(userIDStr string, requesterID string) (int, error)
	UpdateAUser(userData models.UpdateUserRequestModel, userIDStr string, requesterID string) (*models.User, int, error)
	GetAllUsers(pagination database.Pagination) ([]models.User, *database.PaginationResponse, int, error)
	GetUserDataPrivacySettings(userIDStr string, requesterID string) (*models.DataPrivacySettings, int, error)
	UpdateUserDataPrivacySettings(userData models.DataPrivacySettings, userIDStr string, requesterID string) (*models.DataPrivacySettings, int, error)
	UpdateARegion(userData models.UserRegionTimezoneLanguage, userIDStr string, requesterID string) (*models.UserRegionTimezoneLanguage, int, error)
	GetUserRegion(userIDStr string, requesterID string) (*models.UserRegionTimezoneLanguage, int, error)
	ReplaceUserRole(userID string, roleID int) (gin.H, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}
