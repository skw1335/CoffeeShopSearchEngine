package mailer

const (
  FromName = "CoffeeMap"

)
type Client interface {
  Send(templateFile, username, email string, data any, isTestenv bool) error
}
