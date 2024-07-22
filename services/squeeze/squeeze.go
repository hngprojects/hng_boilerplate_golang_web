package service

import (
    "errors"
    "github.com/go-playground/validator/v10"
    "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
    "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
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
    if err := storage.DB.Postgresql.Create(&userSubmission).Error; err != nil {
        return errors.New("failed to save submission to database")
    }
    return nil
}


func SendConfirmationEmail(email string) error {
    subject := "Confirmation Email"
    body := "Thank you for your submission. We have received your request and will process it shortly."
    err := utility.SendEmail(email, subject, body)
    if err != nil {
        return errors.New("failed to send confirmation email")
    }
    return nil
}