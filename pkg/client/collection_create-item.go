package client

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

/*
	CreateItem ( IN Dict<String,Variant> properties,
	             IN Secret secret,
	             IN Boolean replace,
	             OUT ObjectPath item,
	             OUT ObjectPath prompt);
*/

// CreateItem creates an Item in a collection
// item = secret + lookup attributes + label
func (collection *Collection) CreateItem(properties map[string]dbus.Variant,
	secretApi *SecretApi, replace bool) (*Item, string, error) {

	client := collection.Parent
	call, err := client.Call("org.freedesktop.secrets", collection.ObjectPath,
		"org.freedesktop.Secret.Collection", "CreateItem", properties, secretApi, replace)

	if err != nil {
		return nil, "", errors.New("dbus call failed. Error: " + err.Error())
	}

	var itemPath, promptPath dbus.ObjectPath

	err = call.Store(&itemPath, &promptPath)

	if err != nil {
		return nil, "",
			errors.New("Type conversion failed in 'CreateItem'. Error: " + err.Error())
	}

	item := NewItem(collection)

	if label, ok := properties["org.freedesktop.Secret.Item.Label"]; ok {
		if label, ok := label.Value().(string); ok {
			item.Label = label
		}
	} else { // No Label
		item.Label = ""
	}

	if attributes, ok := properties["org.freedesktop.Secret.Item.Attributes"]; ok {
		if value, ok := attributes.Value().(map[string]string); ok {
			item.LookupAttributesMutex.Lock()
			item.LookupAttributes = value
			item.LookupAttributesMutex.Unlock()
		} else {
			return nil, "",
				fmt.Errorf("'Attributes' in 'CreateItem' are not 'map[string]string'. Error: %T", attributes.Value())
		}
	}
	item.ObjectPath = itemPath
	item.Created, _ = item.PropertyCreated()
	item.Modified, _ = item.PropertyModified()
	item.Secret.SecretApi = secretApi
	err = collection.AddItem(item)

	if err != nil {
		return nil, "", errors.New("CreateItem failed. Error: " + err.Error())
	}

	return item, string(promptPath), nil
}
