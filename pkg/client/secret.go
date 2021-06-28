package client

// NewSecret returns a new instance of Secret
func NewSecret(parent *Item) *Secret {

	secret := &Secret{}
	secret.Parent = parent
	secret.SecretApi = &SecretApi{}

	return secret
}

// NewSecret returns a new instance of SecretApi
// SecretApi is the exact secret structure accordinf to API
// Secret is a wrapper around SecretApi to hold extra information
func NewSecretApi() *SecretApi {

	return &SecretApi{}
}
