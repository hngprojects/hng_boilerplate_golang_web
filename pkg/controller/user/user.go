package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) GetUser(c *gin.Context) {

	// var (
	// 	req = models.CreateUserRequestModel{}
	// )

	// err := c.ShouldBind(&req)
	// if err != nil {
	// 	rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
	// 	c.JSON(http.StatusBadRequest, rd)
	// 	return
	// }

	// err = base.Validator.Struct(&req)
	// if err != nil {
	// 	rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
	// 	c.JSON(http.StatusUnprocessableEntity, rd)
	// 	return
	// }

	// reqData, err := auth.ValidateCreateUserRequest(req, base.Db.Postgresql)
	// if err != nil {
	// 	rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
	// 	c.JSON(http.StatusBadRequest, rd)
	// 	return
	// }
	// _ = reqData

	// // perform get user request here

	// base.Logger.Info("user created successfully")
	// rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully",)

	// c.JSON(http.StatusOK, rd)
}
