package service

func NewSecretService(parent *Service) *SecretService {
	secretservice := &SecretService{}
	secretservice.Parent = parent
	return secretservice
}
