// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

// create and initialize a new session
func NewItem(parent *Collection) *Item {
	item := &Item{}
	item.Parent = parent
	item.Locked = false
	item.SaveData = parent.SaveData
	item.Secret = NewSecret(item)
	item.LockMutex = new(sync.Mutex)
	item.PropertiesMutex = new(sync.RWMutex)
	item.LookupAttributesMutex = new(sync.RWMutex)
	item.Properties = make(map[string]dbus.Variant)
	item.LookupAttributes = make(map[string]string)

	item.DataMutex = new(sync.RWMutex)

	return item
}

// GetLookupAttribute returns an attribute with given key otherwise null
func (i *Item) GetLookupAttribute(key string) string {
	i.LookupAttributesMutex.Lock()
	defer i.LookupAttributesMutex.Unlock()
	return i.LookupAttributes[key]
}

// CreateMethodFromPath returns a.b.c.Foo when
// item path is /a/b/c/xyz and passed method is 'Foo'
func (i *Item) CreateMethodFromPath(method string) string {
	_, child := Path2Name(string(i.ObjectPath), method)
	return child
}

////////////////////////////// Signals //////////////////////////////

/*
	ItemCreated (OUT ObjectPath item);
*/

func (item *Item) SignalItemCreated() {

	item.Parent.Parent.Connection.Emit("/org/freedesktop/secrets",
		"org.freedesktop.Secret.Collection.ItemCreated",
		item.ObjectPath)

	log.Infof("Emitted 'ItemCreated' signal for item: %v", item.ObjectPath)
}

/*
	ItemDeleted (OUT ObjectPath item);
*/

func (item *Item) SignalItemDeleted() {

	// This is correct. 'org.freedesktop.Secret.Item' doesn't have any signal itself.
	// when an item is deleted, signal would be emitted from 'Secret.Collection'
	item.Parent.Parent.Connection.Emit("/org/freedesktop/secrets",
		"org.freedesktop.Secret.Collection.ItemDeleted",
		item.ObjectPath)

	log.Infof("Emitted 'ItemDeleted' signal for item: %v", item.ObjectPath)
}

/*
	ItemChanged (OUT ObjectPath item);
*/

func (item *Item) SignalItemChanged() {

	item.Parent.Parent.Connection.Emit("/org/freedesktop/secrets",
		"org.freedesktop.Secret.Collection.ItemChanged",
		item.ObjectPath)

	log.Infof("Emitted 'ItemChanged' signal for item: %v", item.ObjectPath)
}

////////////////////////////// Properties //////////////////////////////

// GetProperty returns given dbus property value
func (item *Item) GetProperty(name string) (dbus.Variant, error) {

	busObject := item.Parent.Parent.Connection.Object("org.freedesktop.secrets",
		item.ObjectPath)

	variant, err := busObject.GetProperty("org.freedesktop.Secret.Item." + name)

	if err != nil {
		return dbus.MakeVariant(nil), fmt.Errorf("error getting property '%s'. Error: %v", name, err)
	}

	item.Properties[name] = variant

	return variant, nil
}

// SetProperty sets given dbus property name to given value
func (item *Item) SetProperty(name string, value interface{}) {

	item.DbusProperties.SetMust("org.freedesktop.Secret.Item", name, value)

	item.DataMutex.Lock()
	item.Properties[name] = dbus.MakeVariant(value)
	item.DataMutex.Unlock()
}

// SetProperties processes raw properties and sets collection.Properties
func (item *Item) SetProperties(properties map[string]dbus.Variant) {

	const base string = "org.freedesktop.Secret.Item."
	processedProperties := make(map[string]dbus.Variant)

	item.rawProperties = properties // keep a copy of original properties

	attributes, ok := properties["org.freedesktop.Secret.Item.Attributes"]

	item.LookupAttributesMutex.Lock()
	if !ok {
		item.LookupAttributes = map[string]string{}
	} else {
		if attributes, ok := attributes.Value().(map[string]string); ok {
			result := map[string]string{}
			for k, v := range attributes {
				result[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
			item.LookupAttributes = result
		} else {
			item.LookupAttributes = map[string]string{}
		}
	}
	item.LookupAttributesMutex.Unlock()

	for k, v := range properties {
		if len(k) < 28 || k[:28] != base {
			continue
		}
		key := strings.TrimSpace(k[28:])
		if len(key) == 0 ||
			key == "Locked" ||
			key == "Created" ||
			key == "Modified" ||
			key == "Attributes" {
			continue
		}
		processedProperties[key] = v
	}

	if label, ok := processedProperties["Label"]; ok {
		if label, ok := label.Value().(string); ok {
			processedProperties["Label"] = dbus.MakeVariant(strings.TrimSpace(label))
		} else { // Label is not string!
			processedProperties["Label"] = dbus.MakeVariant("")
		}
	} else { // Create empty Label
		processedProperties["Label"] = dbus.MakeVariant("")
	}

	item.Properties = processedProperties
	item.Label = item.Properties["Label"].Value().(string)

}

// UpdateModified updated 'Modified' dbus property of this collection
func (item *Item) UpdateModified() {

	item.DataMutex.Lock()
	item.Modified = uint64(time.Now().Unix())
	item.DataMutex.Unlock()

	item.DbusProperties.SetMust("org.freedesktop.Secret.Item",
		"Modified", item.Modified)
}

// Lock locks a collection and updates dbus 'Locked' and 'Modified' properties
func (item *Item) Lock() {
	item.LockMutex.Lock()
	defer item.LockMutex.Unlock()
	item.Locked = true
	item.SetProperty("Locked", true)
}

// Unlock unlocks a collection and updates dbus 'Locked' and 'Modified' properties
func (item *Item) Unlock() {
	item.LockMutex.Lock()
	defer item.LockMutex.Unlock()
	item.Locked = false
	item.SetProperty("Locked", false)
}
