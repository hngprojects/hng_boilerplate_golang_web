package send_invites

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func MockSendInvite(domain, apiKey, senderEmail string, to []string, subject string) error {
	// Initialize Mailgun client
	mg := mailgun.NewMailgun(domain, apiKey)

	// Generate invite link
	inviteLink := MockGenerateInviteLink()

	// Compose the message
	body := fmt.Sprintf("Please click the following link to join: %s", inviteLink)

	message := mg.NewMessage(
		senderEmail,
		subject,
		body,
		to...,
	)

	// Send the email
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	fmt.Printf("Email sent successfully: ID: %s, Response: %s\n", id, resp)
	return nil
}

func MockGenerateInviteLink() string {
	return "http://localhost:8080/invite/1234"
}
