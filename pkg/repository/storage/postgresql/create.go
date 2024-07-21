package postgresql

import (
	"fmt"
    "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"gorm.io/gorm"
)


func CreateUserSubmission(db *gorm.DB, submission *models.UserSubmission) error {
    return db.Create(submission).Error
}


func CreateOneRecord(db *gorm.DB, model interface{}) error {
	result := db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("record creation failed")
	}
	return nil
}

func CreateMultipleRecords(db *gorm.DB, model interface{}, length int) error {
	result := db.Create(model)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != int64(length) {
		return fmt.Errorf("record creation failed")
	}
	return nil
}
