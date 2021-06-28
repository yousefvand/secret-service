package client

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

/*
	Delete ( OUT ObjectPath prompt);
*/

// Delete removes the collection
func (collection *Collection) Delete() (dbus.ObjectPath, error) {

	call, err := collection.Parent.Call("org.freedesktop.secrets", collection.ObjectPath,
		"org.freedesktop.Secret.Collection", "Delete")

	if err != nil {
		return "/", errors.New("dbus call failed. Error: " + err.Error())
	}

	var prompt dbus.ObjectPath

	err = call.Store(&prompt)

	if err != nil {
		return "/",
			errors.New("Type conversion failed in 'Delete' collection. Error: " + err.Error())
	}

	client := collection.Parent
	err = client.RemoveCollection(collection)

	if err != nil {
		return "/", fmt.Errorf("cannot remove collection '%v'. Error: %v",
			string(collection.ObjectPath), err.Error())
	}

	return prompt, nil
}
