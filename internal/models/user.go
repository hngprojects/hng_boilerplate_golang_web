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
	Role          int            `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
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

type ChangePasswordRequestModel struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=7"`
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

	query := db.Where("id = ?", userID)
	query = postgresql.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

	if err := query.First(&user).Error; err != nil {
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

func (u *User) GetSeedUsers(db *gorm.DB) ([]User, error) {
	var users []User

	query := postgresql.PreloadEntities(db, &users, "Profile", "Products", "Organisations")
	query = query.Limit(2)

	if err := query.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (u *User) Update(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &u)
	return err
}
