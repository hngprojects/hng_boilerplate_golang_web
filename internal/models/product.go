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
	Image       string     `gorm:"column:image; type:text" json:"image"`
	Category    []Category `gorm:"many2many:product_categories;;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	CreatedAt   time.Time  `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type CreateProductRequestModel struct {
	Image       string  `json:"image"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Category    string  `json:"category" validate:"required"`
}

type DeleteProductRequestModel struct {
	ProductID string `json:"product_id" validate:"required"`
}
type UpdateProductRequestModel struct {
	Image       string  `json:"image"`
	ProductID   string  `json:"product_id" validate:"required"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required"`
}

type FilterProduct struct {
	Price    float64 `json:"price"`
	Category string  `json:"category"`
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

func (p *Product) DeleteProduct(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) GetProduct(db *gorm.DB, id string) (Product, error) {
	var product Product
	err := db.Preload("Category").Model(p).First(&product, "id = ?", id).Error
	if err != nil {
		return Product{}, err
	}

	return product, nil
}
