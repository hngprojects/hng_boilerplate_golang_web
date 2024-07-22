package models

import "time"

type UserSubmission struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    Email         string    `gorm:"type:varchar(100);uniqueIndex" json:"email" binding:"required,email"`
    FirstName     string    `json:"first_name" binding:"required"`
    LastName      string    `json:"last_name" binding:"required"`
    Phone         string    `json:"phone" binding:"required"`
    Location      string    `json:"location" binding:"required"`
    JobTitle      string    `json:"job_title" binding:"required"`
    Company       string    `json:"company" binding:"required"`
    Interests     []string  `gorm:"type:text[]" json:"interests" binding:"required"`
    ReferralSource string   `json:"referral_source" binding:"required"`
    CreatedAt     time.Time `json:"created_at"`
}