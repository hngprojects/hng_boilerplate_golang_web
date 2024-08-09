package test_organisation

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func initialise(currUUID string, t *testing.T, r *gin.Engine, db *storage.Database, user auth.Controller, org organisation.Controller, status bool) (string, string) {
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	tst.SignupUser(t, r, user, userSignUpData, status)

	token := tst.GetLoginToken(t, r, user, loginData)

	organisationCreationData := models.CreateOrgRequestModel{
		Name:        fmt.Sprintf("Org %v", currUUID),
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		Description: "Some random description about vibranium",
		State:       "test",
		Industry:    "user",
		Type:        "type1",
		Address:     "wakanda land",
		Country:     "wakanda",
	}

	orgID := tst.CreateOrganisation(t, r, db, org, organisationCreationData, token)

	return orgID, token
}

func SetupOrgTestRouter() (*gin.Engine, *organisation.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tst.Setup()
	db := storage.Connection()
	validator := validator.New()

	orgController := &organisation.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupOrgRoutes(r, orgController)
	return r, orgController
}

func SetupOrgRoutes(r *gin.Engine, orgController *organisation.Controller) {
	orgUrl := r.Group("/api/v1",
		middleware.Authorize(orgController.Db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))

	orgUrl.POST("/organisations/:org_id/roles", orgController.CreateOrgRole)
	orgUrl.GET("/organisations/:org_id/roles", orgController.GetOrgRoles)
	orgUrl.GET("/organisations/:org_id/roles/:role_id", orgController.GetAOrgRole)
	orgUrl.DELETE("/organisations/:org_id/roles/:role_id", orgController.DeleteOrgRole)
	orgUrl.PATCH("/organisations/:org_id/roles/:role_id", orgController.UpdateOrgRole)
	orgUrl.PATCH("/organisations/:org_id/roles/:role_id/permissions", orgController.UpdateOrgPermissions)
}
