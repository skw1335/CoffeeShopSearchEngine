package mailer

import (
  "github.com/sendgrid/sendgrid-go"
)
type SendGridMailer struct { 
  fromEmail string
  apiKey    string
  client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
  client := sendgrid.NewSendClient(apiKey)

  return &SendGridMailer{
    fromEmail: fromEmail,
    apiKey:    apiKey,
    client:    client,
  }
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isTestenv bool) error {
  from :=  mail.NewEmail(FromName, m.fromEmail)
  to := mail.NewEmail(username, email)

  //template parsing and building 
  subject := new(bytes.Buffer)
  body := new(bytes.Buffer)

  message := mail.NewSingleEmail(from, subject, to, "", body)
  
  // 
