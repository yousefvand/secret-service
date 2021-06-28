package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	SearchItems ( IN Dict<String,String> attributes,
	              OUT Array<ObjectPath> results);
*/

// SearchItems Searches for items in this collection matching the lookup attributes
func (collection *Collection) SearchItems(attributes map[string]string) ([]dbus.ObjectPath, error) {

	call, err := collection.Parent.Call("org.freedesktop.secrets", collection.ObjectPath,
		"org.freedesktop.Secret.Collection", "SearchItems", attributes)

	if err != nil {
		return nil, errors.New("dbus call failed. Error: " + err.Error())
	}

	var results []dbus.ObjectPath

	err = call.Store(&results)

	if err != nil {
		return nil,
			errors.New("Type conversion failed in 'SearchItems' collection. Error: " + err.Error())
	}

	return results, nil
}
