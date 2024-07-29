package models

import (
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

func (r *Role) CreateRole(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &r)

	if err != nil {
		return err
	}

	return nil
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
