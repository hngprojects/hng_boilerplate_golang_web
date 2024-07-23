package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Auth(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	auth := auth.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	userUrl := r.Group(fmt.Sprintf("%v/auth", ApiVersion))
	{
		userUrl.POST("/users/signup", auth.CreateUser)
		userUrl.POST("/admin/signup", auth.CreateAdmin)
		userUrl.POST("/login", auth.LoginUser)
		userUrl.POST("/password-reset", auth.ResetPassword)
		userUrl.POST("/password-reset/verify", auth.VerifyResetToken)
		userUrl.POST("/change-password", auth.ChangePassword)
		userUrl.POST("/magick-link", auth.RequestMagicLink)
		userUrl.POST("/magick-link/verify", auth.VerifyMagicLink)
	}
	return r
}
