package models

import (
	"errors"
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type UserRegionTimezoneLanguage struct {
	ID         string         `gorm:"type:uuid;primary_key" json:"id"`
	UserID     string         `gorm:"type:uuid;not null" json:"user_id"`
	RegionID   string         `gorm:"type:uuid;not null" json:"region_id" validate:"required"`
	TimezoneID string         `gorm:"type:uuid;not null" json:"timezone_id" validate:"required"`
	LanguageID string         `gorm:"type:uuid;not null" json:"language_id" validate:"required"`
	CreatedAt  time.Time      `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type Language struct {
	ID        string         `gorm:"type:uuid;primary_key" json:"language_id"`
	Name      string         `gorm:"type:varchar(40);unique;not null" json:"name" validate:"required"`
	Code      string         `gorm:"type:varchar(20);unique;not null" json:"code" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Timezone struct {
	ID          string         `gorm:"type:uuid;primary_key" json:"timezone_id"`
	Timezone    string         `gorm:"type:varchar(40);unique;null" json:"timezone" validate:"required"`
	GmtOffset   string         `gorm:"type:varchar(20);unique;null" json:"gmt_offset" validate:"required"`
	Description string         `gorm:"type:varchar(100);null" json:"description" validate:"required"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Region struct {
	ID        string         `gorm:"type:uuid;primary_key" json:"region_id"`
	Name      string         `gorm:"type:varchar(40);unique;not null" json:"name" validate:"required"`
	Code      string         `gorm:"type:varchar(20);unique;not null" json:"code" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (l *Language) BeforeCreate(tx *gorm.DB) (err error) {

	if l.ID == "" {
		l.ID = utility.GenerateUUID()
	}
	return
}

func (t *Timezone) BeforeCreate(tx *gorm.DB) (err error) {

	if t.ID == "" {
		t.ID = utility.GenerateUUID()
	}
	return
}

func (r *Region) BeforeCreate(tx *gorm.DB) (err error) {

	if r.ID == "" {
		r.ID = utility.GenerateUUID()
	}
	return
}

func (u *UserRegionTimezoneLanguage) GetUserRegionByID(db database.DatabaseManager, userID string) (UserRegionTimezoneLanguage, error) {
	var user UserRegionTimezoneLanguage

	query := db.DB().Where("user_id = ?", userID)
	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil

}

func (u *UserRegionTimezoneLanguage) CreateUserRegion(db database.DatabaseManager) error {
	err := db.CreateOneRecord(&u)

	if err != nil {
		return err
	}

	return nil
}

func (l *Language) CreateLanguage(db database.DatabaseManager) error {
	err := db.CreateOneRecord(&l)

	if err != nil {
		return err
	}

	return nil
}

func (t *Timezone) CreateTimeZone(db database.DatabaseManager) error {
	err := db.CreateOneRecord(&t)

	if err != nil {
		return err
	}

	return nil
}

func (r *Region) CreateRegion(db database.DatabaseManager) error {
	err := db.CreateOneRecord(&r)

	if err != nil {
		return err
	}

	return nil
}

func (r *Region) GetRegions(db database.DatabaseManager) ([]Region, error) {
	var regions []Region
	err := db.SelectAllFromDb("desc", "", &regions, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return regions, err
		}
		return regions, err
	}

	return regions, nil
}

func (r *Timezone) GetTimeZones(db database.DatabaseManager) ([]Timezone, error) {
	var timezones []Timezone
	err := db.SelectAllFromDb("desc", "", &timezones, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return timezones, err
		}
		return timezones, err
	}

	return timezones, nil
}

func (r *Language) GetLanguages(db database.DatabaseManager) ([]Language, error) {
	var languages []Language
	err := db.SelectAllFromDb("desc", "", &languages, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return languages, err
		}
		return languages, err
	}

	return languages, nil
}

func (u *UserRegionTimezoneLanguage) UpdateUserRegion(db database.DatabaseManager) error {
	_, err := db.SaveAllFields(&u)
	return err
}

func (t *Timezone) GetTimezoneByID(db database.DatabaseManager, ID string) (Timezone, error) {
	var timezone Timezone

	query := db.DB().Where("id = ?", ID)
	if err := query.First(&timezone).Error; err != nil {
		return timezone, err
	}

	return timezone, nil

}

func (t *Timezone) UpdateTimeZone(db database.DatabaseManager) error {
	_, err := db.SaveAllFields(&t)
	return err
}
