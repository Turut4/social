package mailer

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	subject, body, err := parseTemplate(templateFile, data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject, to, "", body)

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	if err := retry(func() error {
		return m.sendEmail(email, message)
	}, MaxRetries); err != nil {
		return err
	}

	return nil
}

func (m *SendGridMailer) sendEmail(email string, message *mail.SGMailV3) error {
	res, err := m.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email to %v: %w", email, err)
	}

	log.Printf("email sent to %v with status code %v", email, res.StatusCode)
	return nil
}
