package test_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func SetupAuthTestRouter() (*gin.Engine, *auth.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tests.Setup()
	db := storage.Connection()
	validator := validator.New()

	authController := &auth.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupAuthRoutes(r, authController)
	return r, authController
}

func SetupAuthRoutes(r *gin.Engine, userController *auth.Controller) {
	r.PUT("/api/v1/auth/change-password",
		middleware.Authorize(userController.Db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User),
		userController.ChangePassword)
}
