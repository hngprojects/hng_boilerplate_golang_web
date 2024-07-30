package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type Language struct {
	ID        string         `gorm:"type:uuid;primary_key" json:"language_id"`
	Name      string         `gorm:"type:varchar(40);unique;not null" json:"name" validate:"required"`
	Code      string         `gorm:"type:varchar(10);unique;not null" json:"code" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Timezone struct {
	ID         string         `gorm:"type:uuid;primary_key" json:"timezone_id"`
	Identifier string         `gorm:"type:varchar(40);unique;not null" json:"identifier" validate:"required"`
	Offset     string         `gorm:"type:varchar(10);unique;not null" json:"offset" validate:"required"`
	CreatedAt  time.Time      `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type Region struct {
	ID        string         `gorm:"type:uuid;primary_key" json:"region_id"`
	Name      string         `gorm:"type:varchar(40);unique;not null" json:"name" validate:"required"`
	Code      string         `gorm:"type:varchar(10);unique;not null" json:"code" validate:"required"`
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

func (l *Language) CreateLanguage(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &l)

	if err != nil {
		return err
	}

	return nil
}

func (t *Timezone) CreateTimeZone(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)

	if err != nil {
		return err
	}

	return nil
}

func (r *Region) CreateRegion(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &r)

	if err != nil {
		return err
	}

	return nil
}
