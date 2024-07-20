package waitlist

import (
	"net/http"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func SignupWaitlistUserService(req models.CreateWaitlistUserRequest) (*models.WaitlistUser, int, error) {
	user := &models.WaitlistUser{
		ID:    "",
		Name:  req.Name,
		Email: req.Email,
	}

	err := user.CreateWaitlistUser(nil)
	if err != nil {
		code := http.StatusInternalServerError
		if err == models.ErrWaitlistUserExist {
			code = http.StatusBadRequest
		}
		return nil, code, err
	}

	//@TODO: implement email sending her

	return user, http.StatusCreated, nil
}
