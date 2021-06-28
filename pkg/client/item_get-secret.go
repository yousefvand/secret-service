package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	GetSecret ( IN ObjectPath session,
	            OUT Secret secret);
*/

// GetSecret retrieves the secret for this item
func (item *Item) GetSecret(session dbus.ObjectPath) (*SecretApi, error) {

	client := item.Parent.Parent
	call, err := client.Call("org.freedesktop.secrets", item.ObjectPath,
		"org.freedesktop.Secret.Item", "GetSecret", session)

	if err != nil {
		return nil, errors.New("dbus call failed. Error: " + err.Error())
	}

	var secretApi SecretApi

	err = call.Store(&secretApi)

	if err != nil {
		return nil,
			errors.New("Type conversion failed in 'GetSecret' item. Error: " + err.Error())
	}

	return &secretApi, nil
}
