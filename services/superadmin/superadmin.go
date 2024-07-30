package superadmin

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"gorm.io/gorm"
)

func AddToRegion(region *models.Region, db *gorm.DB) error {

	if err := region.CreateRegion(db); err != nil {
		return err
	}

	return nil
}

func AddToTimeZone(timezone *models.Timezone, db *gorm.DB) error {

	if err := timezone.CreateTimeZone(db); err != nil {
		return err
	}

	return nil
}

func AddToLanguage(language *models.Language, db *gorm.DB) error {

	if err := language.CreateLanguage(db); err != nil {
		return err
	}

	return nil
}
