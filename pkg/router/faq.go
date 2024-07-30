package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/faq"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func FAQ(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	faq := faq.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	faqUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		faqUrl.POST("/faq", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), faq.AddToFaq)
		faqUrl.GET("/faq", faq.GetFaq)
		faqUrl.DELETE("/faq/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), faq.DeleteFaq)
		faqUrl.PUT("/faq/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), faq.UpdateFaq)
	}
	return r
}
