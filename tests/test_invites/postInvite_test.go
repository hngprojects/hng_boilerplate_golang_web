package test_invites

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
// 	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
// )

// func TestPostInvite(t *testing.T) {

// 	setup := InviteSetup(t, true)
// 	requestURI := url.URL{Path: "/api/v1/organisations/send-invite"}

// 	tests := []struct {
// 		Name         string
// 		RequestBody  models.InvitationRequest
// 		ExpectedCode int
// 		Message      string
// 	}{
// 		{
// 			Name: "Successful invite",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  setup.OrgID,
// 				Emails: []string{"micahshallom@gmail.com"},
// 			},
// 			ExpectedCode: http.StatusCreated,
// 			Message:      "Invitation(s) sent successfully",
// 		},
// 		{
// 			Name: "Invalid org_id format",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusUnprocessableEntity,
// 			Message:      "Invalid org_id format",
// 		},
// 		{
// 			Name: "Missing OrgID",
// 			RequestBody: models.InvitationRequest{
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "OrgID is required",
// 		},
// 		{
// 			Name: "Empty Emails Array",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{},
// 			},
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "Emails array cannot be empty",
// 		},
// 		{
// 			Name: "Non-existent OrgID",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "non-existent-org-id",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusNotFound,
// 			Message:      "Organization not found",
// 		},
// 		{
// 			Name: "Non-member User Sending Invite",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusForbidden,
// 			Message:      "User is not a member of the organization",
// 		},
// 		{
// 			Name: "Invalid Email Format",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"invalid-email"},
// 			},
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "Invalid email format",
// 		},
// 		{
// 			Name: "Duplicate Emails",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com", "test@example.com"},
// 			},
// 			ExpectedCode: http.StatusConflict,
// 			Message:      "Duplicate emails found",
// 		},
// 		{
// 			Name: "Exceeding Email Limit",
// 			RequestBody: models.InvitationRequest{
// 				OrgID: "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{
// 					"email1@example.com",
// 					"email2@example.com",
// 					// ... add more emails to exceed the limit
// 				},
// 			},
// 			ExpectedCode: http.StatusBadRequest,
// 			Message:      "Email limit exceeded",
// 		},
// 		{
// 			Name: "Unauthorized User",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusUnauthorized,
// 			Message:      "Unauthorized",
// 		},
// 		{
// 			Name: "Expired Token",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusUnauthorized,
// 			Message:      "Token has expired",
// 		},
// 		{
// 			Name: "Database Error Handling",
// 			RequestBody: models.InvitationRequest{
// 				OrgID:  "0190d5be-e185-72ef-b74a-0c9fce0e2328",
// 				Emails: []string{"test@example.com"},
// 			},
// 			ExpectedCode: http.StatusInternalServerError,
// 			Message:      "Internal server error",
// 		},
// 	}

// 	for _, test := range tests {
// 		r := gin.Default()
// 		r.POST("/api/v1/invite", setup.InviteController.PostInvite)

// 		t.Run(test.Name, func(t *testing.T) {
// 			var b bytes.Buffer
// 			json.NewEncoder(&b).Encode(test.RequestBody)

// 			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", "Bearer "+setup.Token)

// 			rr := httptest.NewRecorder()
// 			r.ServeHTTP(rr, req)

// 			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

// 			data := tst.ParseResponse(rr)

// 			code := int(data["code"].(float64))
// 			tst.AssertStatusCode(t, code, test.ExpectedCode)

// 			if test.Message != "" {
// 				message := data["message"]
// 				if message != nil {
// 					tst.AssertResponseMessage(t, message.(string), test.Message)
// 				} else {
// 					tst.AssertResponseMessage(t, "", test.Message)
// 				}

// 			}

// 		})
// 	}
// }
