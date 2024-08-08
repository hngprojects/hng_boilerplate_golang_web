package test_invites

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func TestGetInvites(t *testing.T) {
	setup := InviteSetup(t, true)
	requestURI := url.URL{Path: "/api/v1/organization/invitations"}

	tests := []struct {
		Name         string
		ExpectedCode int
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "Successful Invitation Retrieval",
			ExpectedCode: http.StatusOK,
			Message:      "Invitations Successfully retrieved",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
		{
			Name:         "Invalid Token Test",
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token is invalid!",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + invalidToken,
			},
		},
		{
			Name:         "Missing Authentication Token",
			ExpectedCode: http.StatusUnauthorized,
			Message:      "Token could not be found!",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		inviteURL := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(setup.DB.Postgresql))
		{
			inviteURL.GET("/organization/invitations", setup.InviteController.GetInvites)
		}

		t.Run(test.Name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, requestURI.String(), nil)
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
