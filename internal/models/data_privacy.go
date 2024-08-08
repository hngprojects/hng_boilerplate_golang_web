package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type DataPrivacySettings struct {
	ID                    string         `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	UserID                string         `gorm:"type:uuid;not null" json:"user_id"`
	ProfileVisibility     bool           `gorm:"default:false" json:"profile_visibility"`
	ShareDataWithPartners bool           `gorm:"default:false" json:"share_data_with_partners"`
	ReceiveEmailUpdates   bool           `gorm:"default:false" json:"receive_email_updates"`
	Enable2FA             bool           `gorm:"default:false" json:"enable_2fa"`
	UseDataEncryption     bool           `gorm:"default:false" json:"use_data_encryption"`
	AllowAnalytics        bool           `gorm:"default:false" json:"allow_analytics"`
	PersonalizedAds       bool           `gorm:"default:false" json:"personalized_ads"`
	CreatedAt             time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *DataPrivacySettings) BeforeCreate(tx *gorm.DB) (err error) {

	if d.ID == "" {
		d.ID = utility.GenerateUUID()
	}
	return
}

func (d *DataPrivacySettings) CreateDataPrivacySettings(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &d)

	if err != nil {
		return err
	}

	return nil
}

func (d *DataPrivacySettings) UpdateDataPrivacySettings(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &d)
	return err
}

func (d *DataPrivacySettings) GetUserDataPrivacySettingsByID(db *gorm.DB, userID string) (DataPrivacySettings, error) {
	var user DataPrivacySettings

	query := db.Where("user_id = ?", userID)
	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil

}
