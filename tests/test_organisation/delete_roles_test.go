package test_organisation

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestDeleteOrgRole(t *testing.T) {
	_, orgController := SetupOrgTestRouter()
	db := orgController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	password, _ := utility.HashPassword("password")

	adminUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Admin User",
		Email:    fmt.Sprintf("admin%v@qa.team", currUUID),
		Password: password,
		Role:     int(models.RoleIdentity.SuperAdmin),
	}

	regularUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Regular User",
		Email:    fmt.Sprintf("user%v@qa.team", currUUID),
		Password: password,
		Role:     int(models.RoleIdentity.User),
	}

	orgID := utility.GenerateUUID()

	org := models.Organisation{
		ID:      orgID,
		Name:    fmt.Sprintf("Org comp%v", currUUID),
		Email:   fmt.Sprintf("orgtest%v@qa.team", currUUID),
		OwnerID: adminUser.ID,
	}

	db.Create(&adminUser)
	db.Create(&regularUser)
	db.Create(&org)

	roleID := utility.GenerateUUID()

	role := models.OrgRole{
		ID:             roleID,
		Name:           fmt.Sprintf("Admin Role-%v", utility.RandomString(5)),
		Description:    "Administrator role",
		OrganisationID: orgID,
	}

	db.Create(&role)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, orgController := SetupOrgTestRouter()
		authController := auth.Controller{
			Db:        orgController.Db,
			Validator: orgController.Validator,
			Logger:    orgController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Delete Org Role", func(t *testing.T) {
		router, orgController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *orgController, loginData)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/organisations/%s/roles/%s", orgID, roleID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Role deleted successfully")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/organisations/%s/roles/%s", orgID, roleID), nil)
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})

	t.Run("Forbidden Access - Regular User Trying to Delete Org Role", func(t *testing.T) {
		router, orgController := setup()

		loginData := models.LoginRequestModel{
			Email:    regularUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *orgController, loginData)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/organisations/%s/roles/%s", orgID, roleID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusForbidden)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "not organization owner")
	})

	t.Run("Bad Request - Non-Existent Role", func(t *testing.T) {
		router, orgController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *orgController, loginData)

		nonExistentRoleID := utility.GenerateUUID()

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/organisations/%s/roles/%s", orgID, nonExistentRoleID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "record not found")
	})
}
