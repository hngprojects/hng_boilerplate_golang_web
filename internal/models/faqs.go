package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
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

func (f *FAQ) BeforeCreate(db *gorm.DB) (err error) {

	if f.ID == "" {
		f.ID = utility.GenerateUUID()
	}
	return
}

func (f *FAQ) GetFaqById(db database.DatabaseManager, ID string) (FAQ, error) {
	var faq FAQ

	err, nerr := db.SelectOneFromDb(&faq, "id = ?", ID)
	if nerr != nil {
		return faq, err
	}
	return faq, nil
}

func (f *FAQ) CreateFaq(db database.DatabaseManager) error {

	err := db.CreateOneRecord(&f)

	if err != nil {
		return err
	}

	return nil
}

func (f *FAQ) UpdateFaq(db database.DatabaseManager) error {
	_, err := db.SaveAllFields(&f)
	return err
}

func (f *FAQ) DeleteFaq(db database.DatabaseManager) error {

	err := db.DeleteRecordFromDb(&f)

	if err != nil {
		return err
	}

	return nil
}

func (n *FAQ) FetchAllFaq(db database.DatabaseManager, c *gin.Context) ([]FAQ, database.PaginationResponse, error) {
	var faqs []FAQ

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := db.SelectAllFromDbOrderByPaginated(
		"created_at",
		"desc",
		"",
		pagination,
		&faqs,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return faqs, paginationResponse, nil
}
