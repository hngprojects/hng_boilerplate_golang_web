package models

type Organisation struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Name        string
	Description string
	Users       []*User `gorm:"many2many:org_users;"`
}
