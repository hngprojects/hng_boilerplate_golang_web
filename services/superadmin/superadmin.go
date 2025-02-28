package superadmin

import (
	"net/http"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"gorm.io/gorm"
)

func AddToRegion(region *models.Region, db *gorm.DB) error {

	pdb := inst.InitDB(db)
	if err := region.CreateRegion(pdb); err != nil {
		return err
	}

	return nil
}

func AddToTimeZone(timezone *models.Timezone, db *gorm.DB) error {
	pdb := inst.InitDB(db)

	if err := timezone.CreateTimeZone(pdb); err != nil {
		return err
	}

	return nil
}

func AddToLanguage(language *models.Language, db *gorm.DB) error {

	pdb := inst.InitDB(db)
	if err := language.CreateLanguage(pdb); err != nil {
		return err
	}

	return nil
}

func GetRegions(db *gorm.DB) ([]models.Region, error) {
	pdb := inst.InitDB(db)

	var region models.Region

	regionData, err := region.GetRegions(pdb)
	if err != nil {
		return nil, err
	}

	return regionData, nil
}

func GetTimeZones(db *gorm.DB) ([]models.Timezone, error) {

	var timezone models.Timezone

	pdb := inst.InitDB(db)
	timezoneData, err := timezone.GetTimeZones(pdb)
	if err != nil {
		return nil, err
	}

	return timezoneData, nil
}

func GetLanguages(db *gorm.DB) ([]models.Language, error) {

	var language models.Language

	pdb := inst.InitDB(db)
	languageData, err := language.GetLanguages(pdb)
	if err != nil {
		return nil, err
	}

	return languageData, nil
}

func UpdateATimeZone(req *models.Timezone, reqID string, db *gorm.DB) (*models.Timezone, int, error) {

	var (
		timezone models.Timezone
	)
	pdb := inst.InitDB(db)

	timezone, err := timezone.GetTimezoneByID(pdb, reqID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	timezone.Timezone = req.Timezone
	timezone.GmtOffset = req.GmtOffset
	timezone.Description = req.Description

	if err := timezone.UpdateTimeZone(pdb); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &timezone, http.StatusOK, nil
}
