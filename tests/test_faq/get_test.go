package test_faq

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

func TestGetFaq(t *testing.T) {
	setup := func() (*gin.Engine, *auth.Controller) {
		router, faqController := SetupFAQTestRouter()
		authController := auth.Controller{
			Db:        faqController.Db,
			Validator: faqController.Validator,
			Logger:    faqController.Logger,
		}
		return router, &authController
	}
	_, newsController := SetupFAQTestRouter()
	db := newsController.Db.Postgresql

	faq := models.FAQ{
		ID:        utility.GenerateUUID(),
		Question:  fmt.Sprintf("What is the purpose of this %s FAQ?", utility.RandomString(6)),
		Answer:    "To provide answers to frequently asked questions.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(&faq)

	t.Run("Successful Get FAQ", func(t *testing.T) {
		router, _ := setup()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/faq", nil)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "FAQ retrieved successfully")
	})
}
