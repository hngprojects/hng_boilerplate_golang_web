package test_invites

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func TestPostInvite(t *testing.T) {

	setup := InviteSetup(t, false)
	requestURI := url.URL{Path: "/api/v1/organization/send-invite"}

	tests := []struct {
		Name         string
		RequestBody  models.InvitationRequest
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Successful invitation sent",
			RequestBody: models.InvitationRequest{
				OrgID:  setup.OrgID,
				Emails: []string{setup.Email},
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Invitation(s) sent successfully",
		},
		{
			Name: "Invalid org_id format",
			RequestBody: models.InvitationRequest{
				OrgID:  "0190d5be-e185-72ef-b74a-0",
				Emails: []string{"test@example.com"},
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Validation failed",
		},
		{
			Name: "Missing OrgID",
			RequestBody: models.InvitationRequest{
				Emails: []string{"test@example.com"},
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "Validation failed",
		},
		{
			Name: "Empty Emails Array",
			RequestBody: models.InvitationRequest{
				OrgID:  setup.OrgID,
				Emails: []string{},
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "No emails provided",
		},
	}

	for _, test := range tests {
		r := gin.Default()
		r.POST(requestURI.Path, middleware.Authorize(setup.DB.Postgresql), setup.InviteController.PostInvite)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(test.RequestBody)
			if err != nil {
				t.Fatalf("Failed to encode request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+setup.Token)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)
			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

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