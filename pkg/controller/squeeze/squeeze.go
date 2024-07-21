package controller

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
    "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
    "github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

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

    // Save to database
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
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to submit your request", "status_code": 500})
        return
    }

    err := utility.SendEmail(req.Email, "Subject", "Plain Text Content", "HTML Content")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send email", "status_code": 500})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Your request has been received. You will get a template shortly."})
}