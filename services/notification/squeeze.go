package notifications

import (
	"encoding/json"
	"fmt"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/send"
)

func (n NotificationObject) SendSqueeze() error {
	var (
		notificationData     = models.SendSqueeze{}
		templateFileName     = "squeeze.html"
		baseTemplateFileName = ""
		subject              = "Subject: Welcome to Our Service"
	)

	err := json.Unmarshal([]byte(n.Notification.Data), &notificationData)
	if err != nil {
		return fmt.Errorf("error decoding saved notification data, %v", err)
	}

	data, err := ConvertToMapAndAddExtraData(notificationData, map[string]interface{}{"firstname": thisOrThatStr(notificationData.FirstName, notificationData.Email)})
	if err != nil {
		return fmt.Errorf("error converting data to map, %v", err)
	}

	return send.SendEmail(n.ExtReq, notificationData.Email, subject, templateFileName, baseTemplateFileName, data)
}
