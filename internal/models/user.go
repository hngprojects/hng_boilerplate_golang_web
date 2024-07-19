package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Email    string
	Password string // Consider hashing passwords before storing them
	Profile  Profile
	Products []Product       `gorm:"foreignKey:UserID"`
	Orgs     []*Organisation `gorm:"many2many:org_users;"`
}
