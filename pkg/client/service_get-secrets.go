package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
GetSecrets ( IN Array<ObjectPath> items,
             IN ObjectPath session,
             OUT Dict<ObjectPath,Secret> secrets);
*/

// GetSecrets returns secrets associated to given object paths
func (client *Client) GetSecrets(items []dbus.ObjectPath,
	session dbus.ObjectPath) (map[dbus.ObjectPath]SecretApi, error) {

	call, err := client.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service", "GetSecrets", items, session)

	if err != nil {
		return nil, errors.New("dbus call failed. Error: " + err.Error())
	}

	var secrets map[dbus.ObjectPath]SecretApi

	err = call.Store(&secrets)

	if err != nil {
		return nil, errors.New("Type conversion failed in 'GetSecrets'. Error: " + err.Error())
	}

	return secrets, nil
}
