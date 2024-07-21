package postgresql

import "gorm.io/gorm"

func DeleteRecordFromDb(db *gorm.DB, record interface{}) error {
	tx := db.Delete(record)
	return tx.Error
}

func HardDeleteRecordFromDb(db *gorm.DB, record interface{}) error {
	tx := db.Unscoped().Delete(record)
	return tx.Error
}
