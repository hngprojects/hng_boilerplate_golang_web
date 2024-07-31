package models

type SendOTP struct {
	Email    string `json:"email"  validate:"required"`
	OtpToken int    `json:"otp_token"  validate:"required"`
}

type SendWelcomeMail struct {
	Email int `json:"email"  validate:"required"`
}

