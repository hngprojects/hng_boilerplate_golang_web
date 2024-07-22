package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
    "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
    "github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

var validate = validator.New()

type SqueezeRequest struct {
    Email          string   `json:"email" binding:"required,email"`
    FirstName      string   `json:"first_name" binding:"required"`
    LastName       string   `json:"last_name" binding:"required"`
    Phone          string   `json:"phone" binding:"required"`
    Location       string   `json:"location" binding:"required"`
    JobTitle       string   `json:"job_title" binding:"required"`
    Company        string   `json:"company" binding:"required"`
    Interests      []string `json:"interests" binding:"required"`
    ReferralSource string   `json:"referral_source" binding:"required"`
}

func HandleSqueeze(c *gin.Context) {
    var req SqueezeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to submit your request", "status_code": 400})
        return
    }

    if err := service.ValidateSqueezeRequest(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status_code": 400})
        return
    }

    
    if err := service.CreateUserSubmission(storage.DB.Postgresql, req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to submit your request", "status_code": 500})
        return
    }

    if err := service.SendConfirmationEmail(req.Email); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send email", "status_code": 500})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Your request has been received. You will get a template shortly."})
}
