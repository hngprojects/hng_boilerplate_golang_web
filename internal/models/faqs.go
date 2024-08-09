package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type FAQ struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	Question  string         `gorm:"type:varchar(225);not null" json:"question" validate:"required"`
	Answer    string         `gorm:"type:text;not null" json:"answer" validate:"required"`
	Category  string         `gorm:"type:varchar(30);null" json:"category" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UpdateFAQ struct {
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
	Category string `json:"category" validate:"required"`
}

func (f *FAQ) BeforeCreate(tx *gorm.DB) (err error) {

	if f.ID == "" {
		f.ID = utility.GenerateUUID()
	}
	return
}

func (f *FAQ) GetFaqById(db *gorm.DB, ID string) (FAQ, error) {
	var faq FAQ

	err, nerr := postgresql.SelectOneFromDb(db, &faq, "id = ?", ID)
	if nerr != nil {
		return faq, err
	}
	return faq, nil
}

func (f *FAQ) CreateFaq(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &f)

	if err != nil {
		return err
	}

	return nil
}

func (f *FAQ) UpdateFaq(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &f)
	return err
}

func (f *FAQ) DeleteFaq(db *gorm.DB) error {

	err := postgresql.DeleteRecordFromDb(db, &f)

	if err != nil {
		return err
	}

	return nil
}

func (n *FAQ) FetchAllFaq(db *gorm.DB, c *gin.Context) ([]FAQ, postgresql.PaginationResponse, error) {
	var faqs []FAQ

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := postgresql.SelectAllFromDbOrderByPaginated(
		db,
		"created_at",
		"desc",
		pagination,
		&faqs,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return faqs, paginationResponse, nil
}
