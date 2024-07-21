package models

import (
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	ID          string    `gorm:"type:varchar;primaryKey;unique;not null" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Roles       []Role    `gorm:"many2many:role_permissions;" json:"roles"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;null;autoUpdateTime" json:"updated_at"`
}

type CreatePermissionRequestModel struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description"`
}

func (p *Permission) CreatePermission(db *gorm.DB) error {
	return db.Create(p).Error
}

func (p *Permission) GetPermissionByID(db *gorm.DB, permissionID string) (Permission, error) {
	var permission Permission
	err := db.Where("id = ?", permissionID).First(&permission).Error
	return permission, err
}

func GetPermissionsByIDs(db *gorm.DB, permissionIDs []string) ([]Permission, error) {
	var permissions []Permission
	err := db.Where("id IN ?", permissionIDs).Find(&permissions).Error
	return permissions, err
}
