package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	Delete ( OUT ObjectPath Prompt);
*/

// Delete removes an item from a collection
func (item *Item) Delete() (dbus.ObjectPath, error) {

	client := item.Parent.Parent
	_, err := client.Call("org.freedesktop.secrets", item.ObjectPath,
		"org.freedesktop.Secret.Item", "Delete")

	if err != nil {
		return "", errors.New("dbus call failed. Error: " + err.Error())
	}

	err = item.Parent.RemoveItem(item.ObjectPath)

	if err != nil {
		return "", errors.New("Item delete failed. Error: " + err.Error())
	}

	prompt := dbus.ObjectPath("/")

	return prompt, nil
}
