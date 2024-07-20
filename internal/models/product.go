package models

import "time"

type Product struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"product_id"`
	Name        string    `gorm:"column:name; type:varchar(255); not null" json:"name"`
	Description string    `gorm:"column:description;type:text;" json:"description"`
	OwnerID      string    `gorm:"type:uuid;" json:"owner_id"`
	CreatedAt   time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}
