package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/account"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) GetAccountSettings(c *gin.Context) {
	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	resp, err := account.GetAccountSettings(userId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid recovery email", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("got account successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Got account settings", resp)

	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetSecurityQuestions(c *gin.Context) {
	resp := account.GetSecurityQuestions()

	base.Logger.Info("got security questions successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Security Questions", resp)

	c.JSON(http.StatusOK, rd)
}

func (base *Controller) AddRecoveryEmail(c *gin.Context) {
	var req = models.AddRecoveryEmailRequestModel{}

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

	reqData, err := account.ValidateAddRecoveryEmail(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	err = account.AddRecoveryEmail(reqData, userId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid recovery email", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("recovery email added successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Recovery email successfully added", nil)

	c.JSON(http.StatusOK, rd)
}

func (base *Controller) AddRecoveryPhoneNumber(c *gin.Context) {
	var req = models.AddRecoveryPhoneNumberRequestModel{}

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

	reqData, err := account.ValidateAddRecoveryPhoneNumber(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	err = account.AddRecoveryPhone(reqData, userId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid phone number", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("recovery phone number added successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Recovery phone number successfully added", nil)

	c.JSON(http.StatusOK, rd)
}

func (base *Controller) AddSecurityAnswers(c *gin.Context) {
	var req map[string][]map[string]string

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	reqData, err := account.ValidateAddSecurityQuestions(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Could not submit security questions", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	err = account.AddSecurityAnswers(reqData, userId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Could not submit security questions", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("security answers added successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Security answers submitted successfully", nil)

	c.JSON(http.StatusOK, rd)
}

func (base *Controller) UpdateRecoveryOptions(c *gin.Context) {
	var req models.UpdateRecoveryOptionsRequestModel

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	reqData, err := account.ValidateUpdateRecoveryOptions(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid recovery options", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	err = account.UpdateRecoveryOptions(reqData, userId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid recovery options", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("updated recovery options")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Recovery options updated", nil)

	c.JSON(http.StatusOK, rd)
}

func (base *Controller) DeleteRecoveryOptions(c *gin.Context) {
	var req struct{ Options []string }

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	err = account.DeleteRecoveryOptions(req.Options, userId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Error deleting recovery options", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("deleted recovery options")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Recovery options successfully deleted", nil)

	c.JSON(http.StatusOK, rd)
}
