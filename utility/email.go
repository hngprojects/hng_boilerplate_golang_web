package utility

import (
    "bytes"
    "html/template"
    "log"
    "os"
    "strconv"

    "gopkg.in/gomail.v2"
    "github.com/joho/godotenv"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
    t, err := template.ParseFiles(templateFileName)
    if err != nil {
        return "", err
    }
    buf := new(bytes.Buffer)
    if err = t.Execute(buf, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}

func SendEmail(to, subject, templateFileName string, data interface{}) error {
    m := gomail.NewMessage()

    from := os.Getenv("EMAIL_FROM")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    smtpUser := os.Getenv("SMTP_USER")
    smtpPassword := os.Getenv("SMTP_PASSWORD")

    m.SetHeader("From", from)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)

    body, err := parseTemplate(templateFileName, data)
    if err != nil {
        log.Printf("Could not parse email template: %v", err)
        return err
    }

    m.SetBody("text/html", body)

    port, err := strconv.Atoi(smtpPort)
    if err != nil {
        log.Printf("Invalid SMTP port: %v", smtpPort)
        return err
    }

    d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPassword)

    if err := d.DialAndSend(m); err != nil {
        log.Printf("Could not send email to %q: %v", to, err)
        return err
    }
    
    log.Printf("Email sent to %q successfully", to)
    return nil
}
