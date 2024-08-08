package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/newsletter"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Newsletter(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	newsLetter := newsletter.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	newsLetterUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		newsLetterUrl.POST("/newsletter-subscription", newsLetter.SubscribeNewsLetter)
		newsLetterUrl.GET("/newsletter-subscription",
			middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), newsLetter.GetNewsLetters)
		newsLetterUrl.DELETE("/newsletter-subscription/:id",
			middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), newsLetter.DeleteNewsLetter)
		newsLetterUrl.GET("/newsletter-subscription/deleted",
			middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), newsLetter.GetDeletedNewsLetters)
		newsLetterUrl.PATCH("/newsletter-subscription/restore/:id",
			middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), newsLetter.RestoreDeletedNewsLetter)
	}
	return r
}
