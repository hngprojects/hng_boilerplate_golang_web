package postgresql

import (
	"fmt"
	"gorm.io/gorm"
)

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

func AddUserToOrganisation(db *gorm.DB, orgID, userID string) error {
	// Add user to organisation
	err := db.Exec("INSERT INTO user_organisations (org_id, user_id) VALUES (?, ?)", orgID, userID).Error
	if err != nil {
		return err
	}
	return nil
}


