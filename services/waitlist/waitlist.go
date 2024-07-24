package waitlist

import (
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SignupWaitlistUserService(db *gorm.DB, req models.CreateWaitlistUserRequest) (*models.WaitlistUser, int, error) {
	user := &models.WaitlistUser{
		ID:    utility.GenerateUUID(),
		Name:  req.Name,
		Email: req.Email,
	}

	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)

		existingUser := &models.WaitlistUser{Email: req.Email}
		code, err := existingUser.CheckExistsByEmail(db) // replaced from GetWaitlistUserByEmail to CheckExistsByEmail
		if err != nil {
			return nil, code, models.ErrWaitlistUserExist
		}
	}

	err := user.CreateWaitlistUser(db)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, models.ErrWaitlistUserExist) {
			code = http.StatusBadRequest
		}
		return nil, code, err
	}

	//@TODO: implement email sending her

	return user, http.StatusCreated, nil
}
