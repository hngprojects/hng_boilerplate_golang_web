package waitlist

import (
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func GetWaitLists(c *gin.Context, db *gorm.DB) ([]models.WaitlistUser, *database.PaginationResponse, int, error) {

	var waitList models.WaitlistUser

	pdb := inst.InitDB(db)
	waitLists, paginationResponse, err := waitList.FetchAllWaitList(pdb, c)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return waitLists, nil, http.StatusNoContent, nil
		}
		return waitLists, nil, http.StatusBadRequest, err

	}

	return waitLists, &paginationResponse, http.StatusOK, nil

}

func SignupWaitlistUserService(db *gorm.DB, req models.CreateWaitlistUserRequest) (*models.WaitlistUser, int, error) {
	user := &models.WaitlistUser{
		ID:    utility.GenerateUUID(),
		Name:  req.Name,
		Email: req.Email,
	}

	pdb := inst.InitDB(db)
	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)

		existingUser := &models.WaitlistUser{Email: req.Email}
		code, err := existingUser.CheckExistsByEmail(pdb) // replaced from GetWaitlistUserByEmail to CheckExistsByEmail
		if err != nil {
			return nil, code, models.ErrWaitlistUserExist
		}
	}

	err := user.CreateWaitlistUser(pdb)
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
