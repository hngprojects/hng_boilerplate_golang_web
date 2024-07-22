package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/account"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Account(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	account := account.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	accountRouter := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		accountRouter.POST("/account/add-recovery-email", account.AddRecoveryEmail)
		accountRouter.GET("/account/security-questions", account.GetSecurityQuestions)
		accountRouter.POST("/account/recovery-number", account.AddRecoveryPhoneNumber)
		accountRouter.PUT("/account/update-recovery-options", account.UpdateRecoveryOptions)
		accountRouter.POST("/account/submit-security-answers", account.AddSecurityAnswers)
		accountRouter.DELETE("/account/delete-recovery-options", account.DeleteRecoveryOptions)

		// for testing, not part of the issue (probably should be)
		accountRouter.GET("/account/settings/", account.GetAccountSettings)
	}
	return r
}
