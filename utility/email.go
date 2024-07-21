
package utility

import (
    "fmt"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

const sendGridAPIKey = "your-sendgrid-api-key"

func SendEmail(toEmail, subject, plainTextContent, htmlContent string) error {
    from := mail.NewEmail("Your Name", "your-email@example.com")
    to := mail.NewEmail("Recipient Name", toEmail)
    message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
    client := sendgrid.NewSendClient(sendGridAPIKey)
    response, err := client.Send(message)
    if err != nil {
        return err
    }
    if response.StatusCode >= 400 {
        return fmt.Errorf("failed to send email: %s", response.Body)
    }
    return nil
}
