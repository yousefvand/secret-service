package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	Lock ( IN Array<ObjectPath> objects,
	       OUT Array<ObjectPath> locked,
	       OUT ObjectPath Prompt);
*/

// Lock, locks given objects based on their paths and returns an array of locked object paths
func (client *Client) Lock(
	objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, error) {

	call, err := client.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service", "Lock", objects)

	if err != nil {
		return nil, dbus.ObjectPath("/"), errors.New("dbus call failed. Error: " + err.Error())
	}

	var locked []dbus.ObjectPath
	var prompt dbus.ObjectPath

	err = call.Store(&locked, &prompt)

	if err != nil {
		return nil, dbus.ObjectPath("/"),
			errors.New("Type conversion failed in 'Lock'. Error: " + err.Error())
	}

	for _, objectPath := range locked {
		for _, collection := range client.Collections {
			if !collection.Locked {
				if collection.ObjectPath == objectPath {
					collection.Lock()
				}
			}
			for _, item := range collection.Items {
				if !item.Locked {
					if item.ObjectPath == objectPath {
						item.Lock()
					}
				}
			}
		}
	}

	return locked, prompt, nil
}
