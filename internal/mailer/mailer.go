package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"
)

const (
	FromName            = "Chime"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.templ"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}

func parseTemplate(templateFile string, data any) (subject, body string, err error) {
	templ, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return "", "", err
	}

	subjectBuffer := new(bytes.Buffer)
	err = templ.ExecuteTemplate(subjectBuffer, "subject", data)
	if err != nil {
		return "", "", err
	}

	bodyBuffer := new(bytes.Buffer)
	err = templ.ExecuteTemplate(bodyBuffer, "body", data)
	if err != nil {
		return "", "", err
	}

	return subjectBuffer.String(), bodyBuffer.String(), nil
}

func retry(fn func() error, times int) error {
	if times <= 0 {
		return fmt.Errorf("time should be more than zero")
	}

	var err error
	for i := 0; i <= times; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		time.Sleep(time.Second * time.Duration(i+1))
	}

	return fmt.Errorf("failed after %d attempts: %w", times, err)
}
