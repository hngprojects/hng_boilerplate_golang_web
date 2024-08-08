package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

var ErrWaitlistUserExist = errors.New("waitlist user exists")

type WaitlistUser struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateWaitlistUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (w *WaitlistUser) CreateWaitlistUser(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, w)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrWaitlistUserExist
		}
	}

	return err
}

func (w *WaitlistUser) GetWaitlistUserByEmail(db *gorm.DB) (int, error) {
	err, nerr := postgresql.SelectOneFromDb(db, &w, "email = ?", w.Email)
	if nerr != nil {
		return http.StatusBadRequest, nerr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// added this function to check if waitlist user exists already
func (w *WaitlistUser) CheckExistsByEmail(db *gorm.DB) (int, error) {
	exists := postgresql.CheckExists(db, &w, "email = ?", w.Email)
	if exists {
		return http.StatusBadRequest, errors.New("User exists")
	}

	return http.StatusOK, nil
}

func (n *WaitlistUser) FetchAllWaitList(db *gorm.DB, c *gin.Context) ([]WaitlistUser, postgresql.PaginationResponse, error) {
	var waitLists []WaitlistUser

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := postgresql.SelectAllFromDbOrderByPaginated(
		db,
		"created_at",
		"desc",
		pagination,
		&waitLists,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return waitLists, paginationResponse, nil
}
