package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

type Product struct {
	ID          string     `gorm:"type:uuid;primaryKey" json:"product_id"`
	Name        string     `gorm:"column:name; type:varchar(255); not null" json:"name"`
	Price       float64    `gorm:"column:price; type:decimal(10,2);null" json:"price"`
	Description string     `gorm:"column:description; type:text" json:"description"`
	OwnerID     string     `gorm:"type:uuid;" json:"owner_id"`
	Category    []Category `gorm:"many2many:product_categories;;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	CreatedAt   time.Time  `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type CreateProductRequestModel struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}

func (u *Product) CreateProduct(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &u)
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) AddProductToCategory(db *gorm.DB, categories []interface{}) error {
	// Add product to categories
	err := db.Model(p).Association("Categories").Append(categories...)
	if err != nil {
		return err
	}
	return nil
}
