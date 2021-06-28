package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	CreateCollection ( IN Dict<String,Variant> properties,
	                   IN String alias,
	                   OUT ObjectPath collection,
	                   OUT ObjectPath prompt);
*/

// CreateCollection creates a collection for storing items
// item = secret + lookup attributes + label
func (client *Client) CreateCollection(properties map[string]dbus.Variant,
	alias string) (*Collection, dbus.ObjectPath, error) {

	call, err := client.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service", "CreateCollection", properties, alias)

	if err != nil {
		return nil, "", errors.New("dbus call failed. Error: " + err.Error())
	}

	collection, err := NewCollection(client)

	if err != nil {
		return nil, "/", err
	}

	var collectionObjectPath, prompt dbus.ObjectPath

	err = call.Store(&collectionObjectPath, &prompt)

	if err != nil {
		return nil, "",
			errors.New("Type conversion failed in 'CreateCollection'. Error: " + err.Error())
	}

	collection.Alias = alias
	collection.SetProperties(properties)
	collection.ObjectPath = collectionObjectPath
	client.AddCollection(collection)

	return collection, prompt, nil
}
