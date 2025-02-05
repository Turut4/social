package mailer

import "gopkg.in/gomail.v2"

type MailTrapClient struct {
	apiKey    string
	fromEmail string
}

func NewMailTrap(apiKey, fromEmail string) MailTrapClient {
	return MailTrapClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
	}
}

func (m MailTrapClient) Send(templateFile, email, username string, data any, isSandbox bool) (int, error) {
	subject, body, err := parseTemplate(templateFile, data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("subject", subject)
	message.AddAlternative("text/html", body)

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	if err := retry(func() error {
		return dialer.DialAndSend()
	}, MaxRetries); err != nil {
		return -1, err
	}

	return 200, nil
}
