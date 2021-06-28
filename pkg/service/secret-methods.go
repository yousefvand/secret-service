// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import "sync"

// NewSecret returns a new instance of Secret
func NewSecret(parent *Item) *Secret {

	secret := &Secret{}
	secret.Parent = parent
	secret.SaveData = parent.SaveData
	secret.SecretApi = &SecretApi{}

	secret.DataMutex = new(sync.RWMutex)

	return secret
}
