package test_notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func TestNotificationSettings(t *testing.T) {

	setup := NotifSetup(t, false)

	tests := []struct {
		Name         string
		Method       string
		RequestURI   string
		RequestBody  interface{}
		ExpectedCode int
		Message      string
		ErrorField   string
		ErrorMessage string
		Headers      map[string]string
	}{
		{
			Name:         "Get Notification Settings Action",
			RequestURI:   "/api/v1/settings/notification-settings",
			Method:       http.MethodGet,
			ExpectedCode: http.StatusOK,
			Message:      "Notification settings retrieved successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		}, {
			Name:         "Patch Notifications Settings Action",
			RequestURI:   "/api/v1/settings/notification-settings",
			Method:       http.MethodPatch,
			ExpectedCode: http.StatusOK,
			Message:      "Notification settings updated successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		notificationUrl := r.Group(fmt.Sprintf("%v", "/api/v1/settings/"), middleware.Authorize(setup.DB.Postgresql))
		{
			notificationUrl.GET("/notification-settings", setup.NotificationController.GetNotificationSettings)
			notificationUrl.PATCH("/notification-settings", setup.NotificationController.UpdateNotificationSettings)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(test.Method, test.RequestURI, &b)
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
