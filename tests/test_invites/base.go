package test_invites

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/invite"
	orgController "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type InviteSetupResp struct {
	InviteController *invite.Controller
	Router           *gin.Engine
	Token            string
	OrgID            string
	Email            string
	DB               *storage.Database
}

var invalidToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjAxOTBlMWY0LWYwZDQtNzI4NS1hOWY4LTA3ZmE3ZDA5MjZhNyIsImF1dGhvcmlzZWQiOnRydWUsImV4cCI6MTcyMTk1MDY0NCwicm9sZSI6MSwidXNlcl9pZCI6IjAxOTBlMWYzLWViZTktNzI4NC04MGMzLTEwNjg5NTUzYTQ5NyJ9.Ahrh9l0FJAEEaKIHnph54tdY5U8dEGQiYKiFp6g"

func InviteSetup(t *testing.T, admin bool) InviteSetupResp {
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

	inviteController := &invite.Controller{
		Db:        db,
		Validator: validatorRef,
		Logger:    logger,
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData, admin)
	token := tst.GetLoginToken(t, r, auth, loginData)

	//create an organisation
	orgReq := models.CreateOrgRequestModel{
		Name:        fmt.Sprintf("Org %v", currUUID),
		Description: "This is a test organisation",
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		State:       "Lagos",
		Country:     "Nigeria",
		Industry:    "Tech",
		Type:        "Public",
		Address:     "No 1, Test Street",
	}

	org := orgController.Controller{Db: db, Validator: validatorRef, Logger: logger}
	org_id := tst.CreateOrganisation(t, r, db, org, orgReq, token)

	inviteResp := InviteSetupResp{
		InviteController: inviteController,
		Router:           r,
		Token:            token,
		OrgID:            org_id,
		Email:            email,
		DB:               db,
	}

	return inviteResp
}
