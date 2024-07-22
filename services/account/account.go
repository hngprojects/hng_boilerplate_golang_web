package account

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func ValidateAddRecoveryEmail(req models.AddRecoveryEmailRequestModel) (models.AddRecoveryEmailRequestModel, error) {
	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
		formattedMail, checkBool := utility.EmailValid(req.Email)
		if !checkBool {
			return req, fmt.Errorf("email address is invalid")
		}
		req.Email = formattedMail
	}

	return req, nil
}

func ValidateAddSecurityQuestions(req map[string][]map[string]string) (models.AddSecurityQuesionsRequestModel, error) {
	var data models.AddSecurityQuesionsRequestModel
	answers, ok := req["answers"]
	if !ok {
		return models.AddSecurityQuesionsRequestModel{}, fmt.Errorf("Answers field not provided")
	}

	for _, v := range answers {
		questionOne, ok := v["question_1"]
		if ok {
			answerOne, ok := v["answer_1"]
			if !ok {
				return models.AddSecurityQuesionsRequestModel{}, fmt.Errorf("answer for question one not provided")
			}

			data.QuestionOne = questionOne
			data.AnswerOne = answerOne
		}

		questionTwo, ok := v["question_2"]
		if ok {
			answerTwo, ok := v["answer_2"]
			if !ok {
				return models.AddSecurityQuesionsRequestModel{}, fmt.Errorf("answer for question two not provided")
			}

			data.QuestionTwo = questionTwo
			data.AnswerTwo = answerTwo
		}

		questionThree, ok := v["question_3"]
		if ok {
			answerThree, ok := v["answer_3"]
			if !ok {
				return models.AddSecurityQuesionsRequestModel{}, fmt.Errorf("answer for question three not provided")
			}

			data.QuestionThree = questionThree
			data.AnswerThree = answerThree
		}
	}

	return data, nil
}

func ValidateUpdateRecoveryOptions(req models.UpdateRecoveryOptionsRequestModel) (models.UpdateRecoveryOptionsRequestModel, error) {
	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
		formattedMail, checkBool := utility.EmailValid(req.Email)
		if !checkBool {
			return req, fmt.Errorf("email address is invalid")
		}
		req.Email = formattedMail
	}

	if req.PhoneNumber != "" {
		req.PhoneNumber = strings.ToLower(req.PhoneNumber)
		phone, _ := utility.PhoneValid(req.PhoneNumber)
		req.PhoneNumber = phone
	}

	for _, v := range req.Questions {
		questionOne, ok := v["question_1"]
		if ok {
			answerOne, ok := v["answer_1"]
			if !ok {
				return req, fmt.Errorf("answer for question one not provided")
			}

			req.QuestionOne = questionOne
			req.AnswerOne = answerOne
		}

		questionTwo, ok := v["question_2"]
		if ok {
			answerTwo, ok := v["answer_2"]
			if !ok {
				return req, fmt.Errorf("answer for question two not provided")
			}

			req.QuestionTwo = questionTwo
			req.AnswerTwo = answerTwo
		}

		questionThree, ok := v["question_3"]
		if ok {
			answerThree, ok := v["answer_3"]
			if !ok {
				return req, fmt.Errorf("answer for question three not provided")
			}

			req.QuestionThree = questionThree
			req.AnswerThree = answerThree
		}
	}

	return req, nil
}

func ValidateAddRecoveryPhoneNumber(req models.AddRecoveryPhoneNumberRequestModel) (models.AddRecoveryPhoneNumberRequestModel, error) {
	if req.PhoneNumber != "" {
		req.PhoneNumber = strings.ToLower(req.PhoneNumber)
		phone, _ := utility.PhoneValid(req.PhoneNumber)
		req.PhoneNumber = phone
	}

	return req, nil
}

func GetAccountSettings(userID string, db *gorm.DB) (gin.H, error) {
	user := models.User{}

	account, err := user.GetUserAccountSettings(db, userID)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"email":        account.RecoveryOptions.RecoveryEmail,
		"phone_number": account.RecoveryOptions.RecoveryPhone,
		"answers": map[string]string{
			"question_1": account.RecoveryOptions.QuestionOne,
			"answer_1":   account.RecoveryOptions.AnswerOne,
			"question_2": account.RecoveryOptions.QuestionTwo,
			"answer_2":   account.RecoveryOptions.AnswerTwo,
			"question_3": account.RecoveryOptions.QuestionThree,
			"answer_3":   account.RecoveryOptions.AnswerThree,
		},
	}, nil
}

func GetSecurityQuestions() gin.H {
	return gin.H{
		"answers": map[string]string{
			"question_1": "What is your mother's maiden name?",
			"question_2": "In what city were you born?",
			"question_3": "What is the name of your first pet?",
		},
	}
}

func AddRecoveryEmail(req models.AddRecoveryEmailRequestModel, userID string, db *gorm.DB) error {
	accountSettings := models.AccountSettings{}
	return accountSettings.SetRecoveryEmail(db, userID, req.Email)
}

func AddRecoveryPhone(req models.AddRecoveryPhoneNumberRequestModel, userID string, db *gorm.DB) error {
	accountSettings := models.AccountSettings{}
	return accountSettings.SetRecoveryPhoneNumber(db, userID, req.PhoneNumber)
}

func AddSecurityAnswers(req models.AddSecurityQuesionsRequestModel, userID string, db *gorm.DB) error {
	accountSettings := models.AccountSettings{}
	return accountSettings.SetSecurityQuestions(db, userID, req)
}

func DeleteRecoveryOptions(options []string, userID string, db *gorm.DB) error {
	accountSettings := models.AccountSettings{}

	for _, v := range options {
		switch v {
		case "email":
			err := accountSettings.UnsetRecoveryEmail(db, userID)
			if err != nil {
				return err
			}
		case "phone_number":
			err := accountSettings.UnsetRecoveryPhone(db, userID)
			if err != nil {
				return err
			}
		case "security_questions":
			err := accountSettings.UnsetRecoveryQuestions(db, userID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateRecoveryOptions(options models.UpdateRecoveryOptionsRequestModel, userID string, db *gorm.DB) error {
	accountSettings := models.AccountSettings{}

	if len(options.Email) != 0 {
		err := accountSettings.SetRecoveryEmail(db, userID, options.Email)
		if err != nil {
			return err
		}
	}

	if len(options.PhoneNumber) != 0 {
		err := accountSettings.SetRecoveryPhoneNumber(db, userID, options.PhoneNumber)
		if err != nil {
			return err
		}
	}

	if len(options.Questions) != 0 {
		err := accountSettings.SetSecurityQuestions(db, userID, models.AddSecurityQuesionsRequestModel{
			QuestionOne:   options.QuestionOne,
			QuestionTwo:   options.QuestionTwo,
			QuestionThree: options.QuestionThree,
			AnswerOne:     options.AnswerOne,
			AnswerTwo:     options.AnswerTwo,
			AnswerThree:   options.AnswerThree,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
