package superadmin

import (
	"net/http"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"gorm.io/gorm"
)

// SuperAdminService defines methods for handling regions, timezones, and languages.
type SuperAdminService interface {
	AddToRegion(region *models.Region) error
	AddToTimeZone(timezone *models.Timezone) error
	AddToLanguage(language *models.Language) error
	GetRegions() ([]models.Region, error)
	GetTimeZones() ([]models.Timezone, error)
	GetLanguages() ([]models.Language, error)
	UpdateATimeZone(req *models.Timezone, reqID string) (*models.Timezone, int, error)
}

// superAdminService is the concrete implementation of SuperAdminService.
type superAdminService struct {
	db database.DatabaseManager
}

// NewSuperAdminService initializes the service with a database instance.
func NewSuperAdminService(db *gorm.DB) SuperAdminService {
	return &superAdminService{db: inst.InitDB(db)}
}

func (s *superAdminService) AddToRegion(region *models.Region) error {
	return region.CreateRegion(s.db)
}

func (s *superAdminService) AddToTimeZone(timezone *models.Timezone) error {
	return timezone.CreateTimeZone(s.db)
}

func (s *superAdminService) AddToLanguage(language *models.Language) error {
	return language.CreateLanguage(s.db)
}

func (s *superAdminService) GetRegions() ([]models.Region, error) {
	var region models.Region
	return region.GetRegions(s.db)
}

func (s *superAdminService) GetTimeZones() ([]models.Timezone, error) {
	var timezone models.Timezone
	return timezone.GetTimeZones(s.db)
}

func (s *superAdminService) GetLanguages() ([]models.Language, error) {
	var language models.Language
	return language.GetLanguages(s.db)
}

func (s *superAdminService) UpdateATimeZone(req *models.Timezone, reqID string) (*models.Timezone, int, error) {
	var timezone models.Timezone

	timezone, err := timezone.GetTimezoneByID(s.db, reqID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	timezone.Timezone = req.Timezone
	timezone.GmtOffset = req.GmtOffset
	timezone.Description = req.Description

	if err := timezone.UpdateTimeZone(s.db); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &timezone, http.StatusOK, nil
}
