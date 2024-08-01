package notifications

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/send"
)

func SendWelcomeMail(extReq request.ExternalRequest, req models.SendOTP, db *gorm.DB) error {
	var (
		subject              = "Welcome on board!🎉"
		templateFileName     = "welcome-email.html"
		baseTemplateFileName = ""
		errs                 []string
		user                 models.User
	)

	user, err := user.GetUserByEmail(db, req.Email)
	if err != nil {
		return fmt.Errorf("error getting user with account id %v, %v", req.Email, err)
	}

	data, err := ConvertToMapAndAddExtraData(req, map[string]interface{}{"firstname": thisOrThatStr(user.Profile.FirstName, user.Email)})
	if err != nil {
		return fmt.Errorf("error converting data to map, %v", err)
	}
	err = send.SendEmail(extReq, user.Email, subject, templateFileName, baseTemplateFileName, data)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}
	return nil
}
