package models

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

type Billing struct {
	ID        string         `gorm:"type:uuid;primary_key" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Price     float64        `gorm:"column:price; type:decimal(10,2);null" json:"price"`
	CreatedAt time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateBillingRequest struct {
	Name  string  `json:"title" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

type UpdateBillingRequest struct {
	Name  string  `json:"title"`
	Price float64 `json:"price"`
}

type BillingResponse struct {
	BillingID string    `json:"id"`
	Name      string    `json:"title"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Billing) Create(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &b)

	if err != nil {
		return err
	}

	return nil
}

func (b *Billing) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &b)

	if err != nil {
		return err
	}

	return nil
}

func (b *Billing) GetBillingById(db *gorm.DB, BillingId string) (Billing, error) {
	var Billing Billing
	err, nerr := postgresql.SelectOneFromDb(db, &Billing, "id = ?", BillingId)
	if nerr != nil {
		return Billing, err
	}
	return Billing, nil
}

func (b *Billing) GetAllBillings(db *gorm.DB, c *gin.Context) ([]Billing, postgresql.PaginationResponse, error) {
	var Billing []Billing

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := postgresql.SelectAllFromDbOrderByPaginated(
		db,
		"created_at",
		"desc",
		pagination,
		&Billing,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return Billing, paginationResponse, nil
}

func (b *Billing) UpdateBillingById(db *gorm.DB, req UpdateBillingRequest, BillingId string) (*Billing, error) {
	result, err := postgresql.UpdateFields(db, &b, req, BillingId)

	if err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("failed to update Billing")
	}

	return b, nil
}

func (b *Billing) CheckBillingExists(BillingId string, db *gorm.DB) (Billing, error) {
	Billing, err := b.GetBillingById(db, BillingId)
	if err != nil {
		return Billing, err
	}

	return Billing, nil
}
