package service

import (
    "errors"
    "github.com/go-playground/validator/v10"
    "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
    "github.com/hngprojects/hng_boilerplate_golang_web/utility"
    "gorm.io/gorm"
)

var validate = validator.New()

func ValidateSqueezeRequest(req models.SqueezeRequest) error {
    return validate.Struct(req)
}

func CreateUserSubmission(db *gorm.DB, req models.SqueezeRequest) error {
    userSubmission := models.UserSubmission{
        Email:          req.Email,
        FirstName:      req.FirstName,
        LastName:       req.LastName,
        Phone:          req.Phone,
        Location:       req.Location,
        JobTitle:       req.JobTitle,
        Company:        req.Company,
        Interests:      req.Interests,
        ReferralSource: req.ReferralSource,
    }
    return userSubmission.Create(db)
}

func SendConfirmationEmail(email string) error {
    subject := "Confirmation Email"
    data := struct {
        Email string
    }{
        Email: email,
    }
    err := utility.SendEmail(email, subject, "templates/confirmation_email.html", data)
    if err != nil {
        return errors.New("failed to send confirmation email")
    }
    return nil
}