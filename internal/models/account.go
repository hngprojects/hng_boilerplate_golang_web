package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountSettings struct {
	ID              string                 `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          string                 `gorm:"type:uuid;unique;not null" json:"user_id"`
	RecoveryOptions AccountRecoveryOptions `gorm:"foreignKey:AccountID" json:"recovery_options"`
	CreatedAt       time.Time              `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time              `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type AccountRecoveryOptions struct {
	ID            string    `gorm:"type:uuid;primaryKey;autoincrement" json:"id"`
	AccountID     string    `gorm:"type:uuid;" json:"account_setting_id"`
	RecoveryEmail string    `gorm:"column:recovery_email; type:varchar(255)" json:"email"`
	RecoveryPhone string    `gorm:"column:recovery_phone; type:varchar(255)" json:"phone_number"`
	QuestionOne   string    `gorm:"column:question_one;" json:"question_1"`
	AnswerOne     string    `gorm:"column:answer_one;" json:"answer_1"`
	QuestionTwo   string    `gorm:"column:question_two;" json:"question_2"`
	AnswerTwo     string    `gorm:"column:answer_two;" json:"answer_2"`
	QuestionThree string    `gorm:"column:question_three;" json:"question_3"`
	AnswerThree   string    `gorm:"column:answer_three;" json:"answer_3"`
	CreatedAt     time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type AddRecoveryEmailRequestModel struct {
	Email string `json:"email" validate:"required"`
}

type AddRecoveryPhoneNumberRequestModel struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

// Get account settings gets the user's account settings (obviously).
// I didn't want to modify the already existing GetUserByID method so
// I don't break any existing api
func (u *User) GetUserAccountSettings(db *gorm.DB, userID string) (AccountSettings, error) {
	var accountSettings AccountSettings
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).First(&accountSettings).Error
	if err != nil {
		return AccountSettings{}, err
	}

	return accountSettings, nil
}

// Add Recovery Email as the name implies adds a recovery email to the user's account.
// It returns an error if something absolutely catastrophic happens but finger's crossed
func (a *AccountSettings) AddRecoveryEmail(db *gorm.DB, userID string, email string) error {
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).First(&a).Error
	if err != nil {
		return err
	}

	a.RecoveryOptions.RecoveryEmail = email
	err = db.Save(&a.RecoveryOptions).Error
	if err != nil {
		return err
	}

	return nil
}

// Add Recovery Phone Number is exactly like adding recovery email but with phone number instead of emails.
// It returns an error if something abyssmal happens, let's hope nothing does
func (a *AccountSettings) AddRecoveryPhoneNumber(db *gorm.DB, userID string, phoneNumber string) error {
	err := db.Joins("RecoveryOptions").Where("user_id = ?", userID).First(&a).Error
	if err != nil {
		return err
	}

	a.RecoveryOptions.RecoveryPhone = phoneNumber
	err = db.Save(&a.RecoveryOptions).Error
	if err != nil {
		return err
	}

	return nil
}

type AddSecurityQuesionsParam struct {
	question_one string
	answer_one   string

	question_two string
	answer_two   string

	question_three string
	answer_three   string
}

// Add security questions adds (you guessed it) security questions to the User's account.
// returns an error if something pretty bad happens, but by Carmack's grace it doesn't.
func (a *AccountSettings) AddSecurityQuestions(db *gorm.DB, userID string, questions AddSecurityQuesionsParam) error {
	err := db.Joins("RecoveryOptions").Where("user_id = ?", userID).First(&a).Error
	if err != nil {
		return err
	}

	a.RecoveryOptions.QuestionOne = questions.question_one
	a.RecoveryOptions.AnswerOne = questions.answer_one

	a.RecoveryOptions.QuestionTwo = questions.question_two
	a.RecoveryOptions.AnswerTwo = questions.answer_two

	a.RecoveryOptions.QuestionThree = questions.question_three
	a.RecoveryOptions.AnswerThree = questions.answer_three

	err = db.Save(&a.RecoveryOptions).Error
	if err != nil {
		return err
	}

	return nil
}
