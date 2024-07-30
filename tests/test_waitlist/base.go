package test_waitlist

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/waitlist"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupWLTestRouter() (*gin.Engine, *waitlist.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
	db := storage.Connection()
	validator := validator.New()

	wlController := &waitlist.Controller{
		DB:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupWLRoutes(r, wlController)
	return r, wlController
}

func SetupWLRoutes(r *gin.Engine, wlController *waitlist.Controller) {
	r.GET("/api/v1/waitlist", middleware.Authorize(wlController.DB.Postgresql, models.RoleIdentity.SuperAdmin), wlController.GetWaitLists)
}
