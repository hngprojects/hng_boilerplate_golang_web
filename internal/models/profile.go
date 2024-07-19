package models

type Profile struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"uniqueIndex"`
	FirstName string
	LastName  string
	Phone     int
	AvatarURL string
}
