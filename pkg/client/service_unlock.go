package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	Unlock ( IN Array<ObjectPath> objects,
	         OUT Array<ObjectPath> unlocked,
	         OUT ObjectPath prompt);
*/

// Unlock, unlocks given objects based on their paths and returns an array of unlocked object paths
func (client *Client) Unlock(
	objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, error) {

	call, err := client.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service", "Unlock", objects)

	if err != nil {
		return nil, dbus.ObjectPath("/"), errors.New("dbus call failed. Error: " + err.Error())
	}

	var unlocked []dbus.ObjectPath
	var prompt dbus.ObjectPath

	err = call.Store(&unlocked, &prompt)

	if err != nil {
		return nil, dbus.ObjectPath("/"),
			errors.New("Type conversion failed in 'Unlock'. Error: " + err.Error())
	}

	for _, objectPath := range unlocked {
		for _, collection := range client.Collections {
			if collection.Locked {
				if collection.ObjectPath == objectPath {
					collection.Unlock()
					collection.Modified, _ = collection.PropertyModified()
				}
			}
			for _, item := range collection.Items {
				if item.Locked {
					if item.ObjectPath == objectPath {
						item.Unlock()
						item.Modified, _ = item.PropertyModified()
					}
				}
			}
		}
	}

	return unlocked, prompt, nil
}
