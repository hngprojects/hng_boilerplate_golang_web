package models

import (
	"time"

	"gorm.io/gorm"
)

type Organisation struct {
	ID          string    `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Users       []User    `gorm:"many2many:user_organisations;foreignKey:ID;joinForeignKey:org_id;References:ID;joinReferences:user_id" json:"users"`
	CreatedAt   time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

func (o *Organisation) GetOrgByID(db *gorm.DB, orgID string) (Organisation, error) {
	var org Organisation

	if err := db.Where("id = ?", orgID).First(&org).Error; err != nil {
		return org, err
	}
	return org, nil
}
