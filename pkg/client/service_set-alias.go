package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	SetAlias ( IN String name,
	           IN ObjectPath collection);
*/

// SetAlias sets (or removes) an alias for given collection
func (client *Client) SetAlias(name string, collection dbus.ObjectPath) error {

	_, err := client.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service", "SetAlias", name, collection)

	if err != nil {
		return errors.New("dbus call failed. Error: " + err.Error())
	}

	return nil
}
