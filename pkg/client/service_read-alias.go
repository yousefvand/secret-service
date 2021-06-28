package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	ReadAlias ( IN String name,
	            OUT ObjectPath collection);
*/

// ReadAlias returns the collection with given alias
func (client *Client) ReadAlias(name string) (dbus.ObjectPath, error) {

	call, err := client.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service", "ReadAlias", name)

	if err != nil {
		return dbus.ObjectPath("/"), errors.New("dbus call failed. Error: " + err.Error())
	}

	var collectionPath dbus.ObjectPath

	err = call.Store(&collectionPath)

	if err != nil {
		return dbus.ObjectPath("/"),
			errors.New("Type conversion failed in 'ReadAlias'. Error: " + err.Error())
	}

	return collectionPath, nil
}
