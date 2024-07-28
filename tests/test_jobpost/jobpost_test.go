package test_jobpost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/jobpost"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestJobPostCreate(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/jobs"}
	currUUID := utility.GenerateUUID()
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

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData, false)

	token := tst.GetLoginToken(t, r, auth, loginData)

	tests := []struct {
		Name         string
		RequestBody  models.CreateJobPostModel
		ExpectedCode int
		Message      string
		ErrorField   string
		ErrorMessage string
		Headers      map[string]string
	}{
		{
			Name: "Successful job post creation",
			RequestBody: models.CreateJobPostModel{
				Title:               "Software Engineer Intern",
				Salary:              "5000-7000 USD",
				JobType:             "internship",
				Location:            "San Francisco, CA",
				Deadline:            time.Now().AddDate(0, 1, 0),
				WorkMode:            "remote",
				Experience:          "Entry level (0-2 years)",
				HowToApply:          "Submit your resume and cover letter to hr@company.com",
				JobBenefits:         "Flexible hours, Remote work, Health insurance",
				CompanyName:         "Tech Innovators",
				Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
				KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
				Qualifications:      "Ability to work solo, Bachelor degree",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Job post created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "Invalid job type",
			RequestBody: models.CreateJobPostModel{
				Title:               "Software Engineer Intern",
				Salary:              "5000-7000 USD",
				JobType:             "",
				Location:            "San Francisco, CA",
				Deadline:            time.Now().AddDate(0, 1, 0),
				WorkMode:            "remote",
				Experience:          "Entry level (0-2 years)",
				HowToApply:          "Submit your resume and cover letter to hr@company.com",
				JobBenefits:         "Flexible hours, Remote work, Health insurance",
				CompanyName:         "Tech Innovators",
				Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
				KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
				Qualifications:      "Ability to work solo, Bachelor degree",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
			ErrorField:   "CreateJobPostModel.JobType",
			ErrorMessage: "JobType is a required field",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name: "User unauthorized",
			RequestBody: models.CreateJobPostModel{
				Title:               "Software Engineer Intern",
				Salary:              "5000-7000 USD",
				JobType:             "internship",
				Location:            "San Francisco, CA",
				Deadline:            time.Now().AddDate(0, 1, 0),
				WorkMode:            "remote",
				Experience:          "Entry level (0-2 years)",
				HowToApply:          "Submit your resume and cover letter to hr@company.com",
				JobBenefits:         "Flexible hours, Remote work, Health insurance",
				CompanyName:         "Tech Innovators",
				Description:         "We are looking for a passionate Software Engineer Intern to join our team. You will be working on exciting projects and gain hands-on experience.",
				KeyResponsibilities: "Develop and maintain web applications, Collaborate with the team on various projects, Participate in code reviews",
				Qualifications:      "Ability to work solo, Bachelor degree",
			},
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
		},
	}

	jobPostController := jobpost.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		jobUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
		{
			jobUrl.POST("/jobs", jobPostController.CreateJobPost)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["status_code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

				if test.ErrorField != "" {
					errorData := data["error"].(map[string]interface{})
					errorMessage := errorData[test.ErrorField].(string)
					tst.AssertResponseMessage(t, errorMessage, test.ErrorMessage)
				}
			}
		})
	}
}
