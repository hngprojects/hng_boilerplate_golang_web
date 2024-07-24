package user

import (
	"net/http"

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

// func (base *Controller) GetUserByID(c *gin.Context) {
//     var (
//         userID = c.Param("userId")
//     )

//     user, err := base.UserService.GetUserByID(userID)
//     if err != nil {
//         var rd utility.Response
//         if gorm.IsRecordNotFoundError(err) {
//             rd = utility.BuildErrorResponse(http.StatusNotFound, "error", "User not found", err, nil)
//         } else {
//             rd = utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Internal server error", err, nil)
//         }
//         c.JSON(rd.StatusCode, rd)
//         return
//     }

//     base.Logger.Info("User retrieved successfully")
//     rd := utility.BuildSuccessResponse(http.StatusOK, "User retrieved successfully", user)
//     c.JSON(http.StatusOK, rd)
// }

func (base *Controller) GetUserByID(c *gin.Context) {
    var (
        userID = c.Param("userid")
    )
	

	respData, err := user.GetUserByID(userID, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
    

    base.Logger.Info("User retrieved successfully")
    rd := utility.BuildSuccessResponse(http.StatusOK, "User retrieved successfully", respData)
    c.JSON(http.StatusOK, rd)
}
