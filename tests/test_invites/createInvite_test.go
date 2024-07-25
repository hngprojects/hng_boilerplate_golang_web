package test_invites

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func TestCreateInvite(t *testing.T) {
	setup := InviteSetup(t, false)
	requestURI := url.URL{Path: "/api/v1/invite/create"}

	tests := []struct {
		Name         string
		RequestBody  models.InvitationCreateReq
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name: "Successful Invitation Creation",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: setup.OrgID,
				Email:          setup.Email,
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Invitation created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
		{
			Name: "Invalid Email Format",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: setup.OrgID,
				Email:          "micah",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
		{
			Name: "Empty Email Field",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: setup.OrgID,
				Email:          "",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
		{
			Name: "Missing Organisation ID",
			RequestBody: models.InvitationCreateReq{
				Email: setup.Email,
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Validation failed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
		{
			Name: "Invalid Organisation ID",
			RequestBody: models.InvitationCreateReq{
				OrganisationID: "0190d9a1-e05e-787d-85ee-bd91a61c6da0",
				Email:          setup.Email,
			},
			ExpectedCode: http.StatusNotFound,
			Message:      "Invalid Organisation ID",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		inviteURL := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(setup.DB.Postgresql))
		{
			inviteURL.POST("/invite/create", setup.InviteController.CreateInvite)

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

			}

		})

	}
}