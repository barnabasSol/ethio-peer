package mail

import (
	"bytes"
	broker "ep-mailing-service/internal/broker/rabbitmq"
	"fmt"
	"text/template"

	"resty.dev/v3"
)

type Service struct {
	client      *resty.Client
	apiHost     string
	apiToken    string
	senderEmail string
}

func NewService(apiHost, apiToken, senderEmail string) *Service {
	client := resty.New()
	return &Service{
		client:      client,
		apiHost:     apiHost,
		apiToken:    apiToken,
		senderEmail: senderEmail,
	}
}

func (s *Service) sendEmail(toEmail, subject, textBody, htmlBody string) error {
	payload := map[string]any{
		"from": map[string]string{
			"email": s.senderEmail,
		},
		"to": []map[string]string{
			{"email": toEmail},
		},
		"subject": subject,
		"text":    textBody,
		"html":    htmlBody,
	}

	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+s.apiToken).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(fmt.Sprintf("%s/api/send", s.apiHost))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("mailtrap API error: %s", resp.String())
	}
	return nil
}

func (s *Service) SendOTP(otpPayload broker.OtpPayload) error {
	const tpl = `
		<html>
			<body>
			  <p>Your OTP code is: <b>{{.OTP}}</b></p>
			  <p>Thank you.</p>
			</body>
		</html>
`

	t, err := template.New("otp").Parse(tpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, otpPayload); err != nil {
		return err
	}

	subject := "Your EthioPeer OTP Code"

	textBody := fmt.Sprintf("Hello,\n\nYour OTP code is: %s\n\nThank you.", otpPayload.OTP)
	htmlBody := buf.String()

	return s.sendEmail(otpPayload.Email, subject, textBody, htmlBody)
}

func (s *Service) SendWelcome(welcomePayload broker.WelcomePayload) error {
	subject := "Welcome to Our Service"

	const textTpl = `Hello {{.Name}},
		Welcome to our service! We're glad to have you.
		Best regards.
		`

	textTemplate, err := template.New("welcomeText").Parse(textTpl)
	if err != nil {
		return err
	}
	var textBuf bytes.Buffer
	if err := textTemplate.Execute(&textBuf, welcomePayload); err != nil {
		return err
	}

	const htmlTpl = `
	<html>
		<body>
			  <h1>Welcome {{.Name}}!</h1>
			  <p>We're glad to have you with us.</p>
			  <p>Best regards,</p>
		</body>
	</html>
`
	htmlTemplate, err := template.New("welcomeHTML").Parse(htmlTpl)
	if err != nil {
		return err
	}
	var htmlBuf bytes.Buffer
	if err := htmlTemplate.Execute(&htmlBuf, welcomePayload); err != nil {
		return err
	}

	return s.sendEmail(welcomePayload.Email, subject, textBuf.String(), htmlBuf.String())
}
