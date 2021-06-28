package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	SearchItems ( IN Dict<String,String> attributes,
	              OUT Array<ObjectPath> unlocked,
	              OUT Array<ObjectPath> locked);
*/

// SearchItems searches for items in this collection matching the lookup attributes
func (client *Client) SearchItems(
	attributes map[string]string) ([]dbus.ObjectPath, []dbus.ObjectPath, error) {

	call, err := client.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service", "SearchItems", attributes)

	if err != nil {
		return nil, nil, errors.New("dbus call failed. Error: " + err.Error())
	}

	var unlocked, locked []dbus.ObjectPath

	err = call.Store(&unlocked, &locked)

	if err != nil {
		return nil, nil,
			errors.New("Type conversion failed in 'SearchItems'. Error: " + err.Error())
	}

	return unlocked, locked, nil
}
