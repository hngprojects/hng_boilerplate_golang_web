package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags User
// @Produce json
// @Success 200 {object} utility.Response "Users retrieved successfully"
// @Failure 500 {object} utility.Response "Error retrieving users"
// @Router /users [get]
func (base *Controller) GetAllUsers(c *gin.Context) {

	usersData, paginationResponse, code, err := service.GetAllUsers(c, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Users retrieved successfully", usersData, paginationResponse)
	c.JSON(http.StatusOK, rd)

}

// @Summary Get a user
// @Description Retrieve a single user by ID
// @Tags User
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} utility.Response "User retrieved successfully"
// @Failure 404 {object} utility.Response "User not found"
// @Router /users/{user_id} [get]
func (base *Controller) GetAUser(c *gin.Context) {

	var (
		userID = c.Param("user_id")
	)

	userData, code, err := service.GetAUser(userID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User retrieved successfully", userData)
	c.JSON(http.StatusOK, rd)
}

// @Summary Get user organizations
// @Description Retrieve organizations associated with a user
// @Tags User
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} utility.Response "User organisations retrieved successfully"
// @Failure 404 {object} utility.Response "User not found"
// @Router /users/{user_id}/organisations [get]
func (base *Controller) GetAUserOrganisation(c *gin.Context) {

	var (
		userID = c.Param("user_id")
	)

	userData, code, err := service.GetAUserOrganisation(userID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User organisations retrieved successfully", userData)
	c.JSON(http.StatusOK, rd)
}

// @Summary Delete a user
// @Description Delete a user by ID
// @Tags User
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} utility.Response "User deleted successfully"
// @Failure 404 {object} utility.Response "User not found"
// @Router /users/{user_id} [delete]
func (base *Controller) DeleteAUser(c *gin.Context) {

	var (
		userID = c.Param("user_id")
	)

	code, err := service.DeleteAUser(userID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "User deleted successfully", nil)
	c.JSON(http.StatusOK, rd)
}

// @Summary Update a user
// @Description Update user information
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body models.UpdateUserRequestModel true "User update request"
// @Success 200 {object} utility.Response "User info updated successfully"
// @Failure 400 {object} utility.Response "Failed to parse request body"
// @Failure 422 {object} utility.Response "Validation failed"
// @Router /users/{user_id} [put]
func (base *Controller) UpdateAUser(c *gin.Context) {
	var (
		userID = c.Param("user_id")
		req    = models.UpdateUserRequestModel{}
	)

	err := c.ShouldBind(&req)
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

	respData, code, err := service.UpdateAUser(req, userID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("user info updated successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "User info updated successfully", respData)
	c.JSON(http.StatusOK, rd)

}
