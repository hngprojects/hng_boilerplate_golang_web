package models

import (
	"fmt"
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
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
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
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

type AddSecurityQuesionsRequestModel struct {
	QuestionOne   string `json:"question_1"`
	QuestionTwo   string `json:"question_2"`
	QuestionThree string `json:"question_3"`

	AnswerOne   string `json:"answer_1"`
	AnswerTwo   string `json:"answer_2"`
	AnswerThree string `json:"answer_3"`
}

type UpdateRecoveryOptionsRequestModel struct {
	Email       string              `json:"email"`
	PhoneNumber string              `json:"phone_number"`
	Questions   []map[string]string `json:"security_questions"`

	QuestionOne   string
	QuestionTwo   string
	QuestionThree string
	AnswerOne     string
	AnswerTwo     string
	AnswerThree   string
}

// Get account settings gets the user's account settings (obviously).
// I didn't want to modify the already existing GetUserByID method so
// I don't break any existing api
func (u *User) GetUserAccountSettings(db *gorm.DB, userID string) (AccountSettings, error) {
	var accountSettings AccountSettings
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).Attrs(AccountSettings{ID: utility.GenerateUUID(), UserID: userID}).FirstOrCreate(&accountSettings).Error
	if err != nil {
		return AccountSettings{}, err
	}

	return accountSettings, nil
}

// Add Recovery Email as the name implies adds a recovery email to the user's account.
// It returns an error if something absolutely catastrophic happens but finger's crossed
func (a *AccountSettings) SetRecoveryEmail(db *gorm.DB, userID string, email string) error {
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).Attrs(AccountSettings{ID: utility.GenerateUUID(), UserID: userID}).FirstOrCreate(&a).Error
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	if len(a.RecoveryOptions.ID) == 0 {
		options := AccountRecoveryOptions{
			ID:            utility.GenerateUUID(),
			AccountID:     a.ID,
			RecoveryEmail: email,
		}

		err := db.Create(&options).Error
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}
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
func (a *AccountSettings) SetRecoveryPhoneNumber(db *gorm.DB, userID string, phoneNumber string) error {
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).Attrs(AccountSettings{ID: utility.GenerateUUID(), UserID: userID}).FirstOrCreate(&a).Error
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	if len(a.RecoveryOptions.ID) == 0 {
		options := AccountRecoveryOptions{
			ID:            utility.GenerateUUID(),
			AccountID:     a.ID,
			RecoveryPhone: phoneNumber,
		}

		err := db.Create(&options).Error
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}
	}

	a.RecoveryOptions.RecoveryPhone = phoneNumber
	err = db.Save(&a.RecoveryOptions).Error
	if err != nil {
		return err
	}

	return nil
}

// Add security questions adds (you guessed it) security questions to the User's account.
// returns an error if something pretty bad happens, but by Carmack's grace it doesn't.
func (a *AccountSettings) SetSecurityQuestions(db *gorm.DB, userID string, questions AddSecurityQuesionsRequestModel) error {
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).Attrs(AccountSettings{ID: utility.GenerateUUID(), UserID: userID}).FirstOrCreate(&a).Error
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	if len(a.RecoveryOptions.ID) == 0 {
		options := AccountRecoveryOptions{
			ID:        utility.GenerateUUID(),
			AccountID: a.ID,
		}

		err := db.Create(&options).Error
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}

		a.RecoveryOptions = options
	}

	a.RecoveryOptions.QuestionOne = questions.QuestionOne
	a.RecoveryOptions.AnswerOne = questions.AnswerOne

	a.RecoveryOptions.QuestionTwo = questions.QuestionTwo
	a.RecoveryOptions.AnswerTwo = questions.AnswerTwo

	a.RecoveryOptions.QuestionThree = questions.QuestionThree
	a.RecoveryOptions.AnswerThree = questions.AnswerThree

	err = db.Save(&a.RecoveryOptions).Error
	if err != nil {
		return err
	}

	return nil
}

// unset recovery email
func (a *AccountSettings) UnsetRecoveryEmail(db *gorm.DB, userID string) error {
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).Attrs(AccountSettings{ID: utility.GenerateUUID(), UserID: userID}).FirstOrCreate(&a).Error
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	// just create RecoveryOptions model if it doesn't exist
	if len(a.RecoveryOptions.ID) == 0 {
		options := AccountRecoveryOptions{
			ID:        utility.GenerateUUID(),
			AccountID: a.ID,
		}

		return db.Create(&options).Error
	}

	a.RecoveryOptions.RecoveryEmail = ""
	err = db.Save(&a.RecoveryOptions).Error
	if err != nil {
		return err
	}

	return nil
}

// unset recovery phone
func (a *AccountSettings) UnsetRecoveryPhone(db *gorm.DB, userID string) error {
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).Attrs(AccountSettings{ID: utility.GenerateUUID(), UserID: userID}).FirstOrCreate(&a).Error
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	// just create RecoveryOptions model if it doesn't exist
	if len(a.RecoveryOptions.ID) == 0 {
		options := AccountRecoveryOptions{
			ID:        utility.GenerateUUID(),
			AccountID: a.ID,
		}

		return db.Create(&options).Error
	}

	a.RecoveryOptions.RecoveryPhone = ""
	err = db.Save(&a.RecoveryOptions).Error
	if err != nil {
		return err
	}

	return nil
}

// unset recovery question
func (a *AccountSettings) UnsetRecoveryQuestions(db *gorm.DB, userID string) error {
	err := db.Preload("RecoveryOptions").Where("user_id = ?", userID).Attrs(AccountSettings{ID: utility.GenerateUUID(), UserID: userID}).FirstOrCreate(&a).Error
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	// just create RecoveryOptions model if it doesn't exist
	if len(a.RecoveryOptions.ID) == 0 {
		options := AccountRecoveryOptions{
			ID:        utility.GenerateUUID(),
			AccountID: a.ID,
		}

		return db.Create(&options).Error
	}

	a.RecoveryOptions.QuestionOne = ""
	a.RecoveryOptions.AnswerOne = ""

	a.RecoveryOptions.QuestionTwo = ""
	a.RecoveryOptions.AnswerTwo = ""

	a.RecoveryOptions.QuestionThree = ""
	a.RecoveryOptions.AnswerThree = ""

	err = db.Save(&a.RecoveryOptions).Error
	if err != nil {
		return err
	}

	return nil
}
