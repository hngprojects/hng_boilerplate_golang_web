package test_users

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

func TestGetUserOrganisations(t *testing.T) {
	_, userController := SetupUsersTestRouter()
	db := userController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	password, _ := utility.HashPassword("password")

	// Creating test users and organisations
	adminUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Admin User",
		Email:    fmt.Sprintf("admin%v+1@qa1.team", utility.GenerateUUID()),
		Password: password,
		Role:     int(models.RoleIdentity.SuperAdmin),
	}
	regularUser := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Regular User",
		Email:    fmt.Sprintf("user%v@qa1.team", currUUID),
		Password: password,
		Role:     int(models.RoleIdentity.User),
	}
	org1 := models.Organisation{
		ID:      utility.GenerateUUID(),
		Name:    "Organisation 1",
		Email:   fmt.Sprintf("admorg1%v+2@qa1.team", utility.GenerateUUID()),
		OwnerID: adminUser.ID,
	}
	org2 := models.Organisation{
		ID:      utility.GenerateUUID(),
		Name:    "Organisation 2",
		Email:   fmt.Sprintf("adorg%v+1@qa1.team", utility.GenerateUUID()),
		OwnerID: adminUser.ID,
	}

	db.Create(&adminUser)
	db.Create(&regularUser)
	db.Create(&org1)
	db.Create(&org2)

	userOrg1 := UserOrganisation{
		UserID:         regularUser.ID,
		OrganisationID: org1.ID,
	}
	userOrg2 := UserOrganisation{
		UserID:         regularUser.ID,
		OrganisationID: org2.ID,
	}

	db.Create(&userOrg1)
	db.Create(&userOrg2)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, userController := SetupUsersTestRouter()
		authController := auth.Controller{
			Db:        userController.Db,
			Validator: userController.Validator,
			Logger:    userController.Logger,
		}

		return router, &authController
	}

	t.Run("Successful Get User Organisations", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/organizations", regularUser.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "User organisations retrieved successfully")
	})

	t.Run("User Not Found", func(t *testing.T) {
		router, authController := setup()

		loginData := models.LoginRequestModel{
			Email:    adminUser.Email,
			Password: "password",
		}
		token := tests.GetLoginToken(t, router, *authController, loginData)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/organizations", utility.GenerateUUID()), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "user not found")
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		router, _ := setup()

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/organizations", regularUser.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token is invalid!")
	})
}
