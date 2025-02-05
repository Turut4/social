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

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	subject, body, err := parseTemplate(templateFile, data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject, to, "", body)

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var resStatus int
	if err := retry(func() error {
		resStatus, err = m.sendEmail(email, message)
		return err
	}, MaxRetries); err != nil {
		return -1, err
	}

	return resStatus, err
}

func (m *SendGridMailer) sendEmail(email string, message *mail.SGMailV3) (int, error) {
	res, err := m.client.Send(message)
	if err != nil {
		return res.StatusCode, fmt.Errorf("failed to send email to %v: %w", email, err)
	}

	log.Printf("email sent to %v with status code %v", email, res.StatusCode)
	return res.StatusCode, nil
}
