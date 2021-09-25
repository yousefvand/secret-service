package client

import (
	"errors"

	"github.com/yousefvand/secret-service/pkg/crypto"
)

/*
	SetSecret ( IN Secret secret);
*/

// SetSecret sets the secret for this item
func (item *Item) SetSecret(secretApi *SecretApi) error {

	client := item.Parent.Parent
	_, err := client.Call("org.freedesktop.secrets", item.ObjectPath,
		"org.freedesktop.Secret.Item", "SetSecret", *secretApi)

	if err != nil {
		return errors.New("dbus call failed. Error: " + err.Error())
	}

	item.Secret.SecretApi = secretApi
	session := client.GetSessionByPath(secretApi.Session)
	plainSecret, err := crypto.AesCBCDecrypt(secretApi.Parameters,
		secretApi.Value, session.SymmetricKey)

	if err != nil {
		return errors.New("Decryption error: " + err.Error())
	}

	item.Secret.PlainSecret = string(plainSecret)

	return nil
}
