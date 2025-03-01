package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/profile"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	profileService "github.com/hngprojects/hng_boilerplate_golang_web/services/profile"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Profile(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	profileService := profileService.NewProfileService(db.Postgresql.DB())
	product := profile.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq, ProfileService: profileService}

	profileUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql.DB()))
	{
		profileUrl.PATCH("/profile", product.UpdateProfile)
	}

	return r
}
