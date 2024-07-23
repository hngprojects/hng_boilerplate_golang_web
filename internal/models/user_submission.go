package models

import "time"
import "errors"
import "gorm.io/gorm"
import "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"

type UserSubmission struct {
    ID            string    `gorm:"type:uuid;primaryKey" json:"squeeze_id"`
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



func (us *UserSubmission) Create(db *gorm.DB) error {
    if err := postgresql.Create(us).Error; err != nil {
        return errors.New("failed to save submission to database")
    }
    return nil
}