package mail

import (
	broker "ep-mailing-service/internal/broker/rabbitmq"
	"fmt"

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
	payload := map[string]interface{}{
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
	subject := "Your OTP Code"
	textBody := fmt.Sprintf("Hello,\n\nYour OTP code is: %s\n\nThank you.", otpPayload.OTP)
	htmlBody := fmt.Sprintf(
		"<html><body><h1>Your OTP Code</h1><p>Your OTP code is: <b>%s</b></p><p>Thank you.</p></body></html>",
		otpPayload.OTP,
	)
	return s.sendEmail(otpPayload.Email, subject, textBody, htmlBody)
}

func (s *Service) SendWelcome(welcomePayload broker.WelcomePayload) error {
	subject := "Welcome to Our Service"
	textBody := fmt.Sprintf("Hello %s,\n\nWelcome to our service! We're glad to have you.\n\nBest regards.", welcomePayload.Name)
	htmlBody := fmt.Sprintf("<html><body><h1>Welcome %s!</h1><p>We're glad to have you with us.</p><p>Best regards,</p></body></html>", welcomePayload.Name)
	return s.sendEmail(welcomePayload.Email, subject, textBody, htmlBody)
}
