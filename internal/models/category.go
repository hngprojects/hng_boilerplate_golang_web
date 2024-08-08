package models

import "time"

type Category struct {
	ID        string    `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Products  []Product `gorm:"many2many:product_categories;foreignKey:ID;joinForeignKey:category_id;References:ID;joinReferences:product_id"`
	CreatedAt time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}
