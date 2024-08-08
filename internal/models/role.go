package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

type RoleName string
type RoleId int

type DefaultIdentity struct {
	User       RoleId
	SuperAdmin RoleId
}

var RoleIdentity = DefaultIdentity{
	User:       1,
	SuperAdmin: 2,
}

var (
	UserRoleName  RoleName = "user"
	AdminRoleName RoleName = "admin"
)

type Role struct {
	ID          int            `gorm:"primaryKey;type:int" json:"id"`
	Name        string         `gorm:"unique;not null;type:varchar(20)" json:"name" validate:"required"`
	Description string         `gorm:"unique;not null" json:"description" validate:"required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrgRole struct {
	ID             string         `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Name           string         `gorm:"unique;not null;type:varchar(20)" json:"name" validate:"required"`
	Description    string         `gorm:"not null" json:"description" validate:"required"`
	OrganisationID string         `gorm:"not null" json:"-"`
	Permissions    Permission     `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;" json:"permissions"`
	CreatedAt      time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"-"`
	UpdatedAt      time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Permission struct {
	ID             string         `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	RoleID         string         `gorm:"unique;not null" json:"-"`
	Category       string         `gorm:"not null" json:"category"`
	PermissionList PermissionList `gorm:"type:jsonb" json:"permission_list"`
	CreatedAt      time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"-"`
	UpdatedAt      time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type PermissionList map[string]bool

func (p *PermissionList) Scan(value interface{}) error {
	if b, ok := value.([]byte); ok {
		return json.Unmarshal(b, &p)
	}
	return fmt.Errorf("type assertion to []byte failed")
}

func (p PermissionList) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (r *OrgRole) CreateOrgRole(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &r)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgRole) DeleteOrgRole(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &r)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrgRole) UpdateOrgRole(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &r)
	return err
}

func (rp *Permission) UpdateOrgPermissions(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &rp)
	return err
}

func (r *OrgRole) GetOrgRoles(db *gorm.DB, orgID string) ([]OrgRole, error) {
	var orgRoles []OrgRole

	query := db.Where("organisation_id = ?", orgID)
	err := query.Find(&orgRoles).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return orgRoles, nil
		}
		return orgRoles, err
	}

	return orgRoles, nil
}

func (r *OrgRole) GetAOrgRole(db *gorm.DB, orgID, roleID string) (OrgRole, error) {
	var orgRole OrgRole

	query := db.Where("organisation_id = ?", orgID).Where("id = ?", roleID)
	query = postgresql.PreloadEntities(query, &orgRole, "Permissions")

	err := query.First(&orgRole).Error

	if err != nil {
		return orgRole, err
	}

	return orgRole, nil
}

func (r *Role) UpdateUserRole(db *gorm.DB, userId string, roleId int) (*User, error) {
	var user User

	user, err := user.GetUserByID(db, userId)
	if err != nil {
		return nil, err
	}

	user.Role = roleId

	if _, err := postgresql.SaveAllFields(db, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetRoleName(roleId RoleId) RoleName {
	switch roleId {
	case RoleIdentity.User:
		return UserRoleName
	case RoleIdentity.SuperAdmin:
		return AdminRoleName
	default:
		return "unknown"
	}
}
