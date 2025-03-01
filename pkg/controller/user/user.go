package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db          *storage.Database
	Validator   *validator.Validate
	Logger      *utility.Logger
	ExtReq      request.ExternalRequest
	UserService service.UserService
}

func (base *Controller) GetAllUsers(c *gin.Context) {

	pagination := postgresql.GetPagination(c)
	usersData, paginationResponse, code, err := base.UserService.GetAllUsers(pagination)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Users retrieved successfully", usersData, paginationResponse)
	c.JSON(http.StatusOK, rd)

}
func authHelper(c *gin.Context, db storage.Database) (string, error) {
	userIDFromClaims, err := middleware.GetUserClaims(c, db.Postgresql.DB(), "user_id")
	if err != nil {
		return "", err
	}

	requesterID, ok := userIDFromClaims.(string)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	return requesterID, nil
}

func (base *Controller) GetAUser(c *gin.Context) {

	var (
		userID = c.Param("user_id")
	)

	requesterID, err := authHelper(c, *base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", err.Error(), nil, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	userData, code, err := base.UserService.GetAUser(userID, requesterID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User retrieved successfully", userData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetAUserOrganisation(c *gin.Context) {

	requesterID, err := authHelper(c, *base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", err.Error(), nil, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	userIDStr := c.Param("userID")

	userData, code, err := base.UserService.GetAUserOrganisation(userIDStr, requesterID)

	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User organisations retrieved successfully", userData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) DeleteAUser(c *gin.Context) {

	var (
		userID = c.Param("user_id")
	)
	requesterID, err := authHelper(c, *base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", err.Error(), nil, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	code, err := base.UserService.DeleteAUser(userID, requesterID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User deleted successfully", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) UpdateAUser(c *gin.Context) {
	var (
		userID = c.Param("user_id")
		req    = models.UpdateUserRequestModel{}
	)

	requesterID, err := authHelper(c, *base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", err.Error(), nil, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	err = c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed",
			utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, code, err := base.UserService.UpdateAUser(req, userID, requesterID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("user info updated successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "User info updated successfully", respData)
	c.JSON(http.StatusOK, rd)

}
