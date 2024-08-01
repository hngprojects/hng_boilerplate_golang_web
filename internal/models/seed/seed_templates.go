package seed

import (
	"log"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/send"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

var (
	templateFileNames []string = []string{
		"authorization.html", "email_verification.html", "id_verified.html", "send_otp.html", "welcome_password.html", "password_reset_done_mail.html", "wallet-debited.html", "welcome-email.html", "id_not_verified.html", "password_reset_mail.html", "wallet-funded.html",
	}

	title []string = []string{
		"send_otp",
		"send_invitation",
		"send_password_reset",
		"send_password_reset_done",
		"send_wallet_debited",
		"send_wallet_funded",
		"send_welcome",
		"send_welcome_password",
		"send_default",
		"send_id_verified",
		"send_id_not_verified",
	}

	templates []models.EmailTemplate
)

func SeedTemplates() []models.EmailTemplate {
	for i, fileName := range templateFileNames {
		body, err := send.ParseTemplate(request.ExternalRequest{}, fileName, "", map[string]interface{}{})
		if err != nil {
			log.Fatalf("error parsing template %v, %v", fileName, err)
			continue
		}

		templates = append(templates, models.EmailTemplate{
			ID:   utility.GenerateUUID(),
			Name: title[i],
			Body: body,
		})
	}
	return templates
}
