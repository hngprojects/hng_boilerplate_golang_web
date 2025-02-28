package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

var ErrEmailAlreadySubscribed = errors.New("email already subscribed")

type NewsLetter struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email" validate:"required,email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (n *NewsLetter) BeforeCreate(tx database.DatabaseManager) (err error) {

	if n.ID == "" {
		n.ID = utility.GenerateUUID()
	}
	return
}

func (n *NewsLetter) GetNewsLetterById(db database.DatabaseManager, ID string) (NewsLetter, error) {
	var newsletter NewsLetter

	err, nerr := db.SelectOneFromDb(&newsletter, "id = ?", ID)
	if nerr != nil {
		return newsletter, err
	}
	return newsletter, nil
}

func (n *NewsLetter) GetDeletedNewsLetterById(db database.DatabaseManager, ID string) (NewsLetter, error) {
	var newsletter NewsLetter

	err := db.DB().Unscoped().Where("id = ?", ID).First(&newsletter).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return newsletter, fmt.Errorf("newsletter not found: %w", err)
		}
		return newsletter, fmt.Errorf("failed to retrieve newsletter: %w", err)
	}

	return newsletter, nil
}

func (n *NewsLetter) CreateNewsLetter(db database.DatabaseManager) error {

	err := db.CreateOneRecord(&n)

	if err != nil {
		return err
	}

	return nil
}

func (n *NewsLetter) DeleteNewsLetter(db database.DatabaseManager) error {

	err := db.DeleteRecordFromDb(&n)

	if err != nil {
		return err
	}

	return nil
}

func (n *NewsLetter) UpdateNewsLetter(db database.DatabaseManager) error {
	_, err := db.SaveAllFields(&n)
	return err
}

func (n *NewsLetter) FetchAllNewsLetter(db database.DatabaseManager, c *gin.Context) ([]NewsLetter, database.PaginationResponse, error) {
	var newsLetters []NewsLetter

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := db.SelectAllFromDbOrderByPaginated(
		nil,
		"created_at",
		"desc",
		pagination,
		&newsLetters,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return newsLetters, paginationResponse, nil
}

func (n *NewsLetter) FetchAllDeletedNewsLetter(db database.DatabaseManager, c *gin.Context) ([]NewsLetter, database.PaginationResponse, error) {
	var newsLetters []NewsLetter

	pagination := postgresql.GetPagination(c)

	query := db.DB().Unscoped().Where("deleted_at IS NOT NULL")

	paginationResponse, err := db.SelectAllFromDbOrderByPaginated(
		query,
		"created_at",
		"desc",
		pagination,
		&newsLetters,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return newsLetters, paginationResponse, nil
}
