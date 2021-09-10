package service

func NewSecretService(parent *Service) *SecretService {
	secretservice := &SecretService{}
	secretservice.Session = &SecretServiceCLiSession{}
	secretservice.Parent = parent
	return secretservice
}
