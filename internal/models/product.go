package models

type Product struct {
	ID          string `gorm:"type:uuid;primaryKey" json:"product_id"`
	Name        string `gorm:"column:name; type:varchar(255); not null" json:"name"`
	Description string `gorm:"column:description;type:text;" json:"description"`
	Userid      string `gorm:"type:uuid;" json:"user_id"`
}
