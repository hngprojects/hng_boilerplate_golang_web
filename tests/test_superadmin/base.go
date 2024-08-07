package test_superadmin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/superadmin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupSATestRouter() (*gin.Engine, *superadmin.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
	db := storage.Connection()
	validator := validator.New()

	saController := &superadmin.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupSARoutes(r, saController)
	return r, saController
}

func SetupSARoutes(r *gin.Engine, saController *superadmin.Controller) {
	r.POST("/api/v1/regions", middleware.Authorize(saController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		saController.AddToRegion)
	r.POST("/api/v1/timezones", middleware.Authorize(saController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		saController.AddToTimeZone)
	r.POST("/api/v1/languages", middleware.Authorize(saController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		saController.AddToLanguage)
	r.GET("/api/v1/regions", middleware.Authorize(saController.Db.Postgresql),
		saController.GetRegion)
	r.GET("/api/v1/timezones", middleware.Authorize(saController.Db.Postgresql),
		saController.GetTimeZone)
	r.GET("/api/v1/languages", middleware.Authorize(saController.Db.Postgresql),
		saController.GetLanguage)
	r.PATCH("/api/v1/timezones/:id", middleware.Authorize(saController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		saController.UpdateTimeZone)
}
