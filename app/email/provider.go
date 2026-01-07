package email

type Provider interface {
	Send(to, subject, body string) error
}
