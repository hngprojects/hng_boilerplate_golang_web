package test_notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
)

func TestNotification(t *testing.T) {

	setup := NotifSetup(t, false)

	tests := []struct {
		Name         string
		RequestBody  interface{}
		Method       string
		RequestURI   string
		ExpectedCode int
		Message      string
		ErrorField   string
		ErrorMessage string
		Headers      map[string]string
	}{
		{
			Name: "Notification Creation Action",
			RequestBody: models.NotificationReq{
				Message: "Welcome to HNGi8",
			},
			RequestURI:  "/api/v1/notifications/global",
			Method:      http.MethodPost,
			ExpectedCode: http.StatusCreated,
			Message:      "Notification created successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		}, {
			Name: "Fetch All Notifications Action",
			RequestURI: "/api/v1/notifications/all",
			Method:     http.MethodGet,
			ExpectedCode: http.StatusOK,
			Message:      "Notifications retrieved successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		}, {
			Name: "Fetch Unread Notifications Action",
			RequestURI: "/api/v1/notifications/unread",
			Method:     http.MethodGet,
			ExpectedCode: http.StatusOK,
			Message:      "Unread Notifications retrieved successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},{
			Name: "Update Notification Action",
			RequestBody: models.UpdateNotificationReq{
				IsRead: true,
			},
			RequestURI:  fmt.Sprintf("/api/v1/notifications/%s", setup.NotifID),
			Method:      http.MethodPut,
			ExpectedCode: http.StatusOK,
			Message:      "Notification updated successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + setup.Token,
			},
		},
	}

	for _, test := range tests {
		r := gin.Default()

		notificationUrl := r.Group(fmt.Sprintf("%v", "/api/v1/notifications"), middleware.Authorize(setup.DB.Postgresql))
		{
			notificationUrl.POST("/global", setup.NotificationController.CreateNotification)
			notificationUrl.GET("/all", setup.NotificationController.FetchAllNotifications)
			notificationUrl.GET("/unread", setup.NotificationController.FetchUnReadNotifications)
			notificationUrl.PUT("/:notificationId", setup.NotificationController.UpdateNotification)
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
