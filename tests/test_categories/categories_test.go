package testcategories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/category"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/stretchr/testify/assert"
)

func TestGetCategoryNames(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("emmanueluser%v@qa.team", currUUID),
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

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData, false)

	token := tst.GetLoginToken(t, r, auth, loginData)
	category := category.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r = gin.Default()

	categoryUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
	{
		categoryUrl.GET("/categories", category.GetCategoryNames)
	}

	t.Run("GetCategoryNames", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/categories", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Could not parse response: %v", err)
		}

		assert.Equal(t, "success", response["status"])
		assert.Equal(t, float64(200), response["status_code"])
		assert.Equal(t, "Categories fetched successfully", response["message"])

		assert.Contains(t, response, "data")
		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Invalid response format: 'data' is not an object")
		}

		assert.Contains(t, data, "categories")
		categories, ok := data["categories"].([]interface{})
		if !ok {
			t.Fatalf("Invalid response format: 'categories' is not an array")
		}

		for _, category := range categories {
			categoryMap, ok := category.(map[string]interface{})
			if !ok {
				t.Fatalf("Invalid category format in response")
			}
			_, hasName := categoryMap["Name"]
			assert.True(t, hasName, "Category object should have a 'Name' field")
		}

		assert.Contains(t, data, "totalCount")
		assert.Contains(t, data, "page")
		assert.Contains(t, data, "pageSize")

		assert.Equal(t, float64(9), data["totalCount"])
		assert.Equal(t, float64(1), data["page"])
		assert.Equal(t, float64(10), data["pageSize"])
	})
}
