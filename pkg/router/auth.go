package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Auth(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	auth := auth.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	authUrl := r.Group(fmt.Sprintf("%v/auth", ApiVersion))
	{
		authUrl.POST("/users/signup", auth.CreateUser)
		authUrl.POST("/admin/signup", auth.CreateAdmin)
		authUrl.POST("/login", auth.LoginUser)
		authUrl.POST("/password-reset", auth.ResetPassword)
		authUrl.POST("/password-reset/verify", auth.VerifyResetToken)
		authUrl.POST("/magick-link", auth.RequestMagicLink)
		authUrl.POST("/magick-link/verify", auth.VerifyMagicLink)
	}

	authUrlSec := r.Group(fmt.Sprintf("%v/auth", ApiVersion),
		middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
	{
		authUrlSec.POST("/logout", auth.LogoutUser)
		authUrlSec.PUT("/change-password", auth.ChangePassword)
	}
	return r
}
