package models

import (
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID             string       `gorm:"type:varchar;primaryKey;unique;not null" json:"id"`
	RoleName       string       `gorm:"type:varchar(255);not null" json:"role_name"`
	OrganizationID string       `gorm:"type:varchar;not null" json:"organization_id"`
	Organization   Organisation `gorm:"foreignKey:OrganizationID" json:"-"`
	Permissions    []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt      time.Time    `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"column:updated_at;null;autoUpdateTime" json:"updated_at"`
}

type CreateRoleRequestModel struct {
	RoleName       string   `json:"role_name" validate:"required,min=2,max=255"`
	OrganizationID string   `json:"organization_id" validate:"required"`
	PermissionIDs  []string `json:"permission_ids" validate:"required"`
}

func (r *Role) CreateRole(db *gorm.DB) error {
	return db.Create(r).Error
}

func (r *Role) AddPermissionsToRole(db *gorm.DB, permissionIDs []string) error {
	permissions, err := GetPermissionsByIDs(db, permissionIDs)
	if err != nil {
		return err
	}
	return db.Model(r).Association("Permissions").Append(permissions)
}

func (r *Role) GetRoleByID(db *gorm.DB, roleID string) (Role, error) {
	var role Role
	err := db.Preload("Permissions").Where("id = ?", roleID).First(&role).Error
	return role, err
}
