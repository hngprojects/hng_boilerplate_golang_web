package test_testimonial

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestGetUserTestimonials(t *testing.T) {
	setup := func() (*gin.Engine, *auth.Controller) {
		router, testimonialController := SetupTestimonialTestRouter()
		authController := auth.Controller{
			Db:        testimonialController.Db,
			Validator: testimonialController.Validator,
			Logger:    testimonialController.Logger,
		}
		return router, &authController
	}

	_, testimonialController := SetupTestimonialTestRouter()
	db := testimonialController.Db.Postgresql

	user := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "Test User",
		Email:    fmt.Sprintf("user_%s@qa.team", utility.RandomString(6)),
		Password: "hashedpassword",
		Role:     int(models.RoleIdentity.User),
	}

	result := db.DB().Create(&user)
	if result.Error != nil {
		t.Fatalf(" Failed to create test user: %v", result.Error)
	}

	testimonials := []models.Testimonial{
		{ID: utility.GenerateUUID(), UserID: user.ID, Content: "Testimonial 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: utility.GenerateUUID(), UserID: user.ID, Content: "Testimonial 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, testimonial := range testimonials {
		result := db.DB().Create(&testimonial)
		if result.Error != nil {
			t.Fatalf(" Failed to create testimonial: %v", result.Error)
		}
	}

	t.Run("Successful Get User Testimonials", func(t *testing.T) {
		router, _ := setup()

		url := fmt.Sprintf("/api/v1/testimonials/user/%s", user.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "User testimonials retrieved successfully")
	})
}
