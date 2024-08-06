package test_contact

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/contact"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupContactTestRouter() (*gin.Engine, *contact.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
	db := storage.Connection()
	validator := validator.New()

	contactController := &contact.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupContactRoutes(r, contactController)
	return r, contactController
}

func SetupContactRoutes(r *gin.Engine, contactController *contact.Controller) {
	r.POST("/api/v1/contact", contactController.AddToContactUs)
	r.DELETE("/api/v1/contact/:id", middleware.Authorize(contactController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		contactController.DeleteContactUs)
	r.GET("/api/v1/contact",
		middleware.Authorize(contactController.Db.Postgresql, models.RoleIdentity.SuperAdmin), contactController.GetAllContactUs)
	r.GET("/api/v1/contact/id/:id", middleware.Authorize(contactController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		contactController.GetContactUsById)
	r.GET("/api/v1/contact/email/:email", middleware.Authorize(contactController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		contactController.GetContactUsByEmail)

}
