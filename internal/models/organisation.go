package models

import (
	"time"

	"gorm.io/gorm"
)

type Organisation struct {
	ID          string    `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Description string    `gorm:"type:text" json:"description"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Email       string    `gorm:"type:varchar(255);unique" json:"email"`
	State       string    `gorm:"type:varchar(255)" json:"state"`
	Industry    string    `gorm:"type:varchar(255)" json:"industry"`
	Type        string    `gorm:"type:varchar(255)" json:"type"`
	Address     string    `gorm:"type:varchar(255)" json:"address"`
	Country     string    `gorm:"type:varchar(255)" json:"country"`
	OwnerID     string    `gorm:"type:uuid;" json:"owner_id"`
	Users       []User    `gorm:"many2many:user_organisations;"`
	CreatedAt   time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	Roles       []Role    `gorm:"foreignKey:OrganizationID" json:"roles"`
}

type CreateOrgRequestModel struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description"`
	Email       string `json:"email" validate:"required"`
	State       string `json:"state" validate:"required"`
	Industry    string `json:"industry" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Country     string `json:"country" validate:"required"`
}

func (c *Organisation) CreateOrganisation(db *gorm.DB) error {
	err := db.Create(&c).Error
	if err != nil {
		return err
	}
	return nil
}
