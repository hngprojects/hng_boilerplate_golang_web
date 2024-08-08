package test_testimonial

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestCreateTestimonial(t *testing.T) {
	_, testimonialController := SetupTestimonialTestRouter()
	db := testimonialController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	password, _ := utility.HashPassword("password")

	user := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Regular User",
		Email:    fmt.Sprintf("user%v@qa.team", currUUID),
		Password: password,
		Role:     int(models.RoleIdentity.User),
	}

	db.Create(&user)

	setup := func() (*gin.Engine, *auth.Controller) {
		router, testimonialController := SetupTestimonialTestRouter()
		authController := auth.Controller{
			Db:        testimonialController.Db,
			Validator: testimonialController.Validator,
			Logger:    testimonialController.Logger,
		}

		return router, &authController
	}

	router, authController := setup()

	loginData := models.LoginRequestModel{
		Email:    user.Email,
		Password: "password",
	}
	token := tst.GetLoginToken(t, router, *authController, loginData)

	tests := []struct {
		Name         string
		RequestBody  models.TestimonialReq
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful creation of testimonial",
			RequestBody: models.TestimonialReq{
				Name:    "testtestimonial",
				Content: "testcontent",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "testimonial created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name: "Unauthorized Access",
			RequestBody: models.TestimonialReq{
				Name:    "testtestimonial",
				Content: "testcontent",
			},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "Validation failed",
			RequestBody: models.TestimonialReq{
				Content: "testcontent",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/testimonials", &b)

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			tst.AssertStatusCode(t, resp.Code, test.ExpectedCode)
			response := tst.ParseResponse(resp)
			tst.AssertResponseMessage(t, response["message"].(string), test.Message)

		})
	}
}
