package test_users

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

type UserOrganisation struct {
	UserID         string `gorm:"type:uuid;not null" json:"user_id"`
	OrganisationID string `gorm:"type:uuid;not null" json:"organisation_id"`
}

func SetupUsersTestRouter() (*gin.Engine, *user.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tests.Setup()
	db := storage.Connection()
	validator := validator.New()

	userController := &user.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupUsersRoutes(r, userController)
	return r, userController
}

func SetupUsersRoutes(r *gin.Engine, userController *user.Controller) {
	r.PUT("/api/v1/users/:user_id/roles/:role_id",
		middleware.Authorize(userController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		userController.AssignRoleToUser)
	r.GET("/api/v1/users", middleware.Authorize(userController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		userController.GetAllUsers)
	r.GET("/api/v1/users/:user_id",
		middleware.Authorize(userController.Db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User),
		userController.GetAUser)
	r.DELETE("/api/v1/users/:user_id",
		middleware.Authorize(userController.Db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User),
		userController.DeleteAUser)
	r.PUT("/api/v1/users/:user_id",
		middleware.Authorize(userController.Db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User),
		userController.UpdateAUser)
	r.GET("/api/v1/users/:user_id/organisations",
		middleware.Authorize(userController.Db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User),
		userController.GetAUserOrganisation)
	r.PUT("/api/v1/users/:user_id/regions",
		middleware.Authorize(userController.Db.Postgresql, models.RoleIdentity.SuperAdmin),
		userController.UpdateUserRegion)
}
