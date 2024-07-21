package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

type User struct {
	ID            string         `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Name          string         `gorm:"column:name; type:varchar(255)" json:"name"`
	Email         string         `gorm:"column:email; type:varchar(255)" json:"email"`
	Password      string         `gorm:"column:password; type:text; not null" json:"-"`
	Profile       Profile        `gorm:"foreignKey:Userid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"profile"`
	Organisations []Organisation `gorm:"many2many:user_organisations;;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"organisations" ` // many to many relationship
	Products      []Product      `gorm:"foreignKey:OwnerID" json:"products"`
	CreatedAt     time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type CreateUserRequestModel struct {
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	UserName    string `json:"username" validate:"required"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequestModel struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *User) AddUserToOrganisation(db *gorm.DB, user interface{}, orgs []interface{}) error {

	// Add user to organisation
	err := db.Model(user).Association("Organisations").Append(orgs...)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByID(db *gorm.DB, userID string) (User, error) {
	var user User

	if err := db.Preload("Profile").Preload("Products").Preload("Organisations").Where("id = ?", userID).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) CreateUser(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &u)

	if err != nil {
		return err
	}

	return nil
}
