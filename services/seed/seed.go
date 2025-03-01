package seed

import (
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

type SeedService interface {
	GetUser(userIDStr string) ([]models.User, error)
}

type seedService struct {
	db *gorm.DB
}

func NewSeedService(db *gorm.DB) SeedService {
	return &seedService{db: db}
}

func (s *seedService) GetUser(userIDStr string) ([]models.User, error) {
	var (
		user     models.User
		userResp []models.User
	)

	pdb := inst.InitDB(s.db)

	userResp, err := user.GetSeedUsers(pdb)
	if err != nil {
		return userResp, err
	}

	return userResp, nil
}
