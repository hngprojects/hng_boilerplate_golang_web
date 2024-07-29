package test_invites

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAcceptInvite(t *testing.T) {
	setup := InviteSetup(t, false)

	// Create invitation token for testing
	invitationToken := utility.GenerateUUID()
	invitation := models.Invitation{
		ID:             utility.GenerateUUID(),
		OrganisationID: setup.OrgID,
		Token:          invitationToken,
		UserID:         utility.GenerateUUID(),
		IsValid:        true,
		Email:          setup.Email,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	setup.DB.Postgresql.Create(&invitation)

	tests := []struct {
		Name         string
		Method       string
		URL          string
		RequestBody  interface{}
		ExpectedCode int
		// Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Invitation Acceptance (POST)",
			Method:       http.MethodPost,
			URL:          "/api/v1/invite/accept",
			RequestBody:  models.InvitationAcceptReq{InvitationLink: "http://example.com/invite/" + invitationToken},
			ExpectedCode: http.StatusOK,
			// Message:      "Invitation accepted successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
		{
			Name:         "Invalid Invitation Token (POST)",
			Method:       http.MethodPost,
			URL:          "/api/v1/invite/accept",
			RequestBody:  models.InvitationAcceptReq{InvitationLink: "http://example.com/invite/" + "fnwninvalidtoken132334"},
			ExpectedCode: http.StatusBadRequest,
			// Message:      "Invalid or expired invitation link",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
		{
			Name:         "Invalid Invitation Token (GET)",
			Method:       http.MethodGet,
			URL:          "/api/v1/invite/accept/" + "fnwninvalidtoken132334",
			ExpectedCode: http.StatusBadRequest,
			// Message:      "Invalid or expired invitation link",
			Headers: map[string]string{
				"Authorization": "Bearer " + setup.Token,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			r := gin.Default()

			r.POST("/api/v1/invite/accept", middleware.Authorize(setup.DB.Postgresql), setup.InviteController.PostAcceptInvite)
			r.GET("/api/v1/invite/accept/:t", middleware.Authorize(setup.DB.Postgresql), setup.InviteController.GetAcceptInvite)

			var req *http.Request
			var err error

			if test.Method == http.MethodPost {
				var b bytes.Buffer
				json.NewEncoder(&b).Encode(test.RequestBody)
				req, err = http.NewRequest(test.Method, test.URL, &b)
			} else {
				req, err = http.NewRequest(test.Method, test.URL, nil)
			}

			if err != nil {
				t.Fatal(err)
			}

			for k, v := range test.Headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)
			data := tst.ParseResponse(rr)

			code := int(data["status_code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			// if test.Message != "" {
			// 	message := data["message"]
			// 	if message != nil {
			// 		tst.AssertResponseMessage(t, message.(string), test.Message)
			// 	} else {
			// 		tst.AssertResponseMessage(t, "", test.Message)
			// 	}
			// }
		})
	}
}