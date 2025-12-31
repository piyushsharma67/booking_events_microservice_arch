package service

type MockNotifier struct {
	SentMessages []Notification
}

type Notification struct {
	To      string
	Subject string
	Body    string
}

func (m *MockNotifier) SendNotification(to string, subject string, body string) error {
	m.SentMessages = append(m.SentMessages, Notification{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	return nil
}
