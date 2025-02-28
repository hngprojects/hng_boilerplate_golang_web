package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type ContactUs struct {
	ID        string         `gorm:"type:uuid;primary_key;" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name" binding:"required" validate:"required"`
	Email     string         `gorm:"type:varchar(100);not null;index" json:"email" binding:"required" validate:"required,email"`
	Message   string         `gorm:"type:text;not null" json:"message" binding:"required" validate:"required"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *ContactUs) BeforeCreate(tx *gorm.DB) (err error) {

	if c.ID == "" {
		c.ID = utility.GenerateUUID()
	}
	return
}

func (f *ContactUs) GetContactUsById(db database.DatabaseManager, ID string) (ContactUs, error) {
	var contact ContactUs

	err, nerr := db.SelectOneFromDb(&contact, "id = ?", ID)
	if nerr != nil {
		return contact, err
	}
	return contact, nil
}

func (f *ContactUs) GetContactUsByEmail(db database.DatabaseManager, email string) ([]ContactUs, error) {
	var contacts []ContactUs

	err := db.SelectAllFromDb("", "", &contacts, "email = ?", email)
	if err != nil {
		return contacts, err
	}
	return contacts, nil
}

func (c *ContactUs) CreateContactUs(db database.DatabaseManager) error {

	err := db.CreateOneRecord(&c)

	if err != nil {
		return err
	}

	return nil
}

func (c ContactUs) DeleteContactUs(db database.DatabaseManager) error {

	err := db.DeleteRecordFromDb(&c)

	if err != nil {
		return err
	}

	return nil
}

func (cu *ContactUs) FetchAllContactUs(db database.DatabaseManager, c *gin.Context) ([]ContactUs, database.PaginationResponse, error) {
	var contacts []ContactUs

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := db.SelectAllFromDbOrderByPaginated(
		"created_at",
		"desc",
		"",
		pagination,
		&contacts,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return contacts, paginationResponse, nil
}
