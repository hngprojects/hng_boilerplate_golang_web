package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) CreateUser(c *gin.Context) {

	var (
		req = models.CreateUserRequestModel{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	reqData, err := user.ValidateCreateUserRequest(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, code, err := user.CreateUser(reqData, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("user created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", respData)

	c.JSON(code, rd)
}

func (base *Controller) LoginUser(c *gin.Context) {

	var (
		req = models.LoginRequestModel{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, code, err := user.LoginUser(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("user login successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "user login successfully", respData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetAllCustomers(c *gin.Context) {

	limitStr := c.Query("limit")
	pageStr := c.Query("page")

	if limitStr == "" {
		c.JSON(http.StatusBadRequest, utility.BuildErrorResponse(400, "error", "Missing limit parameter", "Missing limit parameter", nil))
		return
	}
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, utility.BuildErrorResponse(400, "error", "Missing page parameter", "Missing page parameter", nil))
			return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, utility.BuildErrorResponse(400, "error", "Invalid or missing limit parameter", err, nil))
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, utility.BuildErrorResponse(400, "error", "Invalid or missing page parameter", err, nil))
		return
	}
	
	
	respData, totalPages, totalItems, err := user.GetAllCustomers(base.Db.Postgresql, page, limit)
	if err != nil {
		rd := utility.BuildErrorResponse(400, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("All Customers fetched successfully")

	response := gin.H {
		"status_code": 200,
		"current_page": page,
		"total_pages": totalPages,
		"limit": limit,
		"total_items": totalItems,
		"data": respData,
	}

	c.JSON(http.StatusOK, response)
}
