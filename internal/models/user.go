package models

import (
	"errors"
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

func GetUserByID(db *gorm.DB, userID string) (User, error) {
	var user User
	result := db.First(&user, "id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, result.Error
	}
	return user, nil
}
