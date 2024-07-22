package utility

import (
    "fmt"
    "gopkg.in/gomail.v2"
)


func SendEmail(to, subject, body string) error {
    m := gomail.NewMessage()

    m.SetHeader("From", "your-email@example.com")
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-password")

    if err := d.DialAndSend(m); err != nil {
        log.Printf("Could not send email to %q: %v", to, err)
        return err
    }
    
    log.Printf("Email sent to %q successfully", to)
    return nil
}