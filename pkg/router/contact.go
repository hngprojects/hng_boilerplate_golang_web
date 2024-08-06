package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/contact"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Contact(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	contact := contact.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	contactUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		contactUrl.POST("/contact", contact.AddToContactUs)
		contactUrl.GET("/contact", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), contact.GetAllContactUs)
		contactUrl.DELETE("/contact/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), contact.DeleteContactUs)
		contactUrl.GET("/contact/id/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), contact.GetContactUsById)
		contactUrl.GET("/contact/email/:email", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), contact.GetContactUsByEmail)

	}
	return r
}
