package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Name          string         `gorm:"column:name; type:varchar(255)" json:"name"`
	Email         string         `gorm:"column:email; type:varchar(255)" json:"email"`
	Profile       Profile        `gorm:"foreignKey:Userid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"profile"`
	Organisations []Organisation `gorm:"many2many:user_organisations;;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"organisations" ` // many to many relationship
	Products      []Product      `gorm:"foreignKey:OwnerID" json:"products"`
	CreatedAt     time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
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