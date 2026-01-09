package email

type Service struct {
	provider Provider
}

func New(provider Provider) *Service {
	return &Service{provider: provider}
}

// Método genérico
func (s *Service) Send(to, subject, body string) error {
	return s.provider.Send(to, subject, body)
}

// Caso de uso concreto
func (s *Service) SendVerification(to, link string) error {

	return s.Send(
		to,
		"Verificá tu email",
		"Hacé click acá:\n\n"+link,
	)
}
