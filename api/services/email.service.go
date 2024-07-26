package services

import (
	"bytes"
	"crypto/tls"
	"drivers-service/config"
	"drivers-service/data/models"
	"drivers-service/pkg/logging"
	"drivers-service/pkg/tools"
	"fmt"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	logger   logging.Logger
	config   *config.Config
	message  *gomail.Message
	from     string
	smtpUser string
	smtpPass string
	smtpHost string
	smtpPort int
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		logger:   logging.NewLogger(cfg),
		config:   cfg,
		message:  gomail.NewMessage(),
		from:     cfg.SMTP.EmailFrom,
		smtpUser: cfg.SMTP.User,
		smtpPass: cfg.SMTP.Pass,
		smtpHost: cfg.SMTP.Host,
		smtpPort: cfg.SMTP.Port,
	}
}

func (e *EmailService) SendEmail(email string, data *models.EmailData, templateName string) error {

	var body bytes.Buffer

	template, err := tools.ParseTemplateDir("./templates")
	if err != nil {
		e.logger.Fatalf("Could not parse template", err)
	}

	template = template.Lookup(templateName)
	err = template.Execute(&body, &data)
	if err != nil {
		return err
	}

	fmt.Printf("body: %v\n", body.String())
	e.message.SetHeader("From", e.from)
	e.message.SetHeader("To", email)
	e.message.SetHeader("Subject", data.Subject)
	e.message.SetBody("text/html", body.String())
	e.message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dialer := gomail.NewDialer(e.smtpHost, e.smtpPort, e.smtpUser, e.smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.config.SMTP.Security}
	// Send Email
	if err := dialer.DialAndSend(e.message); err != nil {
		return err
	}
	return nil
}
