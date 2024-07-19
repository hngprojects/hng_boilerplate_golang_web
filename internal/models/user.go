package models

import (
	"time"
)

type User struct {
	UserID        string         `gorm:"type:uuid;primaryKey" json:"user_id"`
	Name          string         `gorm:"column:name; type:varchar(255)" json:"name"`
	Email         string         `gorm:"column:email; type:varchar(255)" json:"email"`
	Profile       Profile        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"profile"`
	Organisations []Organisation `gorm:"many2many:user_organisations;" json:"organisations" ` // many to many relationship
	Products      []Product      `gorm:"foreignKey:UserID" json:"products"`
	CreatedAt     time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

// type CreateUserRequestModel struct {
// 	FirstName string `json:"first_name" validate:"required"`
// 	LastName  string `json:"last_name" validate:"required"`
// 	Email     string `json:"email" validate:"required"`
// 	Password  string `json:"password" validate:"required"`
// }