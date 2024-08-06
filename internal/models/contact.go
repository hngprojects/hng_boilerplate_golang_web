package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type ContactUs struct {
	ID        string         `gorm:"type:uuid;primary_key;" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name" binding:"required" validate:"required"`
	Email     string         `gorm:"type:varchar(100);not null;index" json:"email" binding:"required" validate:"required,email"`
	Subject   string         `gorm:"type:varchar(255);not null" json:"subject" binding:"required" validate:"required"`
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

func (f *ContactUs) GetContactUsById(db *gorm.DB, ID string) (ContactUs, error) {
	var contact ContactUs

	err, nerr := postgresql.SelectOneFromDb(db, &contact, "id = ?", ID)
	if nerr != nil {
		return contact, err
	}
	return contact, nil
}

func (f *ContactUs) GetContactUsByEmail(db *gorm.DB, email string) ([]ContactUs, error) {
	var contacts []ContactUs

	err := postgresql.SelectAllFromDb(db, "", &contacts, "email = ?", email)
	if err != nil {
		return contacts, err
	}
	return contacts, nil
}

func (c *ContactUs) CreateContactUs(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &c)

	if err != nil {
		return err
	}

	return nil
}

func (c ContactUs) DeleteContactUs(db *gorm.DB) error {

	err := postgresql.DeleteRecordFromDb(db, &c)

	if err != nil {
		return err
	}

	return nil
}

func (cu *ContactUs) FetchAllContactUs(db *gorm.DB, c *gin.Context) ([]ContactUs, postgresql.PaginationResponse, error) {
	var contacts []ContactUs

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := postgresql.SelectAllFromDbOrderByPaginated(
		db,
		"created_at",
		"desc",
		pagination,
		&contacts,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return contacts, paginationResponse, nil
}
