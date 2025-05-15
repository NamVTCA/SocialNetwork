package sms

import "fmt"

type SMSSender interface {
	Send(to string, body string) error
}

type MockSMSSender struct{}

func (m *MockSMSSender) Send(to, body string) error {
	fmt.Printf("📱 Sending SMS to %s | Body: %s\n", to, body)
	return nil
}

func NewMockSMSSender() SMSSender {
	return &MockSMSSender{}
}
