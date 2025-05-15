package email

import "fmt"

type EmailSender interface {
	Send(to string, subject string, body string) error
}

type MockEmailSender struct{}

func (m *MockEmailSender) Send(to, subject, body string) error {
	fmt.Printf("📧 Sending Email to %s | Subject: %s | Body: %s\n", to, subject, body)
	return nil
}

func NewMockEmailSender() EmailSender {
	return &MockEmailSender{}
}
