package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func setupRoleTestRouter() (*gin.Engine, *user.Controller) {
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
	SetupRolesRoutes(r, userController)
	return r, userController
}

func SetupRolesRoutes(r *gin.Engine, userController *user.Controller) {
	r.PUT("/users/:user_id/roles/:role_id", userController.AssignRoleToUser)
}

func TestE2EUpdateUserRole(t *testing.T) {
	router, userController := setupRoleTestRouter()
	db := userController.Db.Postgresql
	currUUID := utility.GenerateUUID()

	user := models.User{
		Email: fmt.Sprintf("testuser%v@qa.team", currUUID),
	}
	db.Create(&user)

	role := models.Role{
		Name: "admin",
	}
	db.Create(&role)

	t.Run("Successful Update", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%s/roles/%d", user.ID, role.ID), nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Role updated successfully")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%s/roles/%d", user.ID, role.ID), nil)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token could not be found!")
		tests.AssertResponseMessage(t, response["error"].(string), "Unauthorized")
	})

	t.Run("User Not Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%d/roles/%d", 999, role.ID), nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["error"].(string), "invalid user")
	})

	t.Run("Role Not Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%s/roles/%d", user.ID, 999), nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["error"].(string), "invalid role")
	})
}
