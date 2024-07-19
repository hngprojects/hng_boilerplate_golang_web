package models

type Product struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Name        string
	Description string
}
