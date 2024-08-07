package test_squeeze

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/squeeze"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupSqueezeTestRouter() (*gin.Engine, *squeeze.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
	db := storage.Connection()
	validator := validator.New()

	squeezeController := &squeeze.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupSqueezeRoutes(r, squeezeController)
	return r, squeezeController
}

func SetupSqueezeRoutes(r *gin.Engine, squeezeController *squeeze.Controller) {
	r.POST("/api/v1/squeeze", squeezeController.Create)
}
