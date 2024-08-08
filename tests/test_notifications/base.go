package test_notifications

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/notificationCRUD"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type NotificationSetupResp struct {
	NotificationController *notificationCRUD.Controller
	Router                 *gin.Engine
	Token                  string
	NotifID                string
	DB                     *storage.Database
}

func NotifSetup(t *testing.T, admin bool) NotificationSetupResp {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	email := fmt.Sprintf("testuser" + currUUID + "@qa.team")

	userSignUpData := models.CreateUserRequestModel{
		Email:       email,
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

	notifController := &notificationCRUD.Controller{
		Db:        db,
		Validator: validatorRef,
		Logger:    logger,
	}


	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData, admin)
	token := tst.GetLoginToken(t, r, auth, loginData)

	//create an organisation
	notReq := models.NotificationReq{
		Message: "Test Notification",
	}

	not := notificationCRUD.Controller{Db: db, Validator: validatorRef, Logger: logger}
	not_id := tst.CreateNotification(t, r, db, not, notReq, token)

	notifResp := NotificationSetupResp{
		NotificationController: notifController,
		Router:                 r,
		Token:                  token,
		NotifID:                not_id,
		DB:                     db,
	}

	return notifResp
}
