// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

// NewCollection creates and initialize a new collection
func NewCollection(parent *Service) *Collection {
	collection := &Collection{}
	collection.Parent = parent
	collection.Locked = false
	collection.SaveData = parent.SaveData
	collection.LockMutex = new(sync.Mutex)
	collection.ItemsMutex = new(sync.RWMutex)
	collection.Items = make(map[string]*Item)
	collection.rawProperties = make(map[string]dbus.Variant)
	collection.Properties = make(map[string]dbus.Variant)

	collection.DataMutex = new(sync.RWMutex)

	return collection
}

// DefaultCollection create and initialize a new default
// collection at: '/org/freedesktop/secrets/aliases/default'
func DefaultCollection(parent *Service, locked bool, created uint64, modified uint64) {

	collection := NewCollection(parent)
	collection.Locked = false
	collection.Alias = "default"
	path := "/org/freedesktop/secrets/aliases/default"
	collection.ObjectPath = dbus.ObjectPath(path)
	collection.Properties = map[string]dbus.Variant{
		"Label": dbus.MakeVariant("default"),
	}

	collection.Parent.AddCollection(collection, locked, created, modified, true)
	// create deafult collection object on dbus
	dbusDefaultCollection(collection, locked, created, modified)
}

// AddItem adds a new item to collection's items
func (collection *Collection) AddItem(item *Item, replace bool, saveData bool,
	locked bool, created uint64, modified uint64, inPlace bool) error {

	collection.ItemsMutex.Lock()
	if replace {
		for _, collectionItem := range collection.Items {
			// Documentation: If replace is set, then it replaces an item already
			//                present with the same values for the attributes
			// Should we check for same Label ?
			if reflect.DeepEqual(collectionItem.LookupAttributes, item.LookupAttributes) {
				collectionItem = item
			}
		}
	} else {
		collection.Items[string(item.ObjectPath)] = item
	}
	collection.ItemsMutex.Unlock()

	if inPlace {
		item.Secret.PlainSecret = string(item.Secret.PlainSecret)
	} else {
		session := collection.Parent.GetSessionByPath(item.Secret.SecretApi.Session)
		if session == nil {
			log.Warn("Secret session is missing")
			return errors.New("Secret session is missing")
		}

		if session.EncryptionAlgorithm == Plain {
			item.Secret.PlainSecret = string(item.Secret.SecretApi.Value)
		} else {
			iv := item.Secret.SecretApi.Parameters
			secret, err := AesCBCDecrypt(iv, item.Secret.SecretApi.Value, session.SymmetricKey)
			if err != nil {
				log.Errorf("Cannot add item due to decryption error. Error: %v", err)
				return errors.New("Decryption error: " + err.Error())
			}
			item.Secret.PlainSecret = string(secret)
		}
		log.Infof("New item at: %v", item.ObjectPath)
	}

	collection.ItemsMutex.Lock()
	collection.Items[string(item.ObjectPath)] = item
	collection.ItemsMutex.Unlock()

	// add item object to dbus
	dbusAddItem(collection, item, locked, created, modified)

	if saveData {
		collection.SaveData()
	}

	return nil
}

// RemoveItem removes an item from collection's item map
func (collection *Collection) RemoveItem(item *Item) {
	collection.ItemsMutex.Lock()
	_, ok := collection.Items[string(item.ObjectPath)]
	if !ok {
		log.Errorf("Item doesn't exist to be removed: %v",
			item.ObjectPath)
		return
	}
	delete(collection.Items, string(item.ObjectPath))
	collection.ItemsMutex.Unlock()
	epoch := Epoch()
	dbusUpdateItems(collection, false, epoch, epoch)
	log.Infof("Item removed: %v", item.ObjectPath)
	collection.SaveData()
}

// GetItemByPath returns the collection with given dbus object path, otherwise null
func (collection *Collection) GetItemByPath(itemPath dbus.ObjectPath) *Item {
	collection.ItemsMutex.RLock()
	defer collection.ItemsMutex.RUnlock()

	for _, item := range collection.Items {
		if item.ObjectPath == itemPath {
			return item
		}
	}
	return nil
}

// CreateMethodFromPath returns a.b.c.Foo when
// collection path is /a/b/c/xyz and passed method is 'Foo'
func (collection *Collection) CreateMethodFromPath(method string) string {
	_, child := Path2Name(string(collection.ObjectPath), method)
	return child
}

////////////////////////////// Signals //////////////////////////////

/*
	CollectionCreated (OUT ObjectPath collection);
*/

// SignalCollectionCreated emits a signal that a new collection was created
func (collection *Collection) SignalCollectionCreated() {

	collection.Parent.Connection.Emit("/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service.CollectionCreated",
		collection.ObjectPath)

	log.Infof("Emitted 'CollectionCreated' signal for collection: %v", collection.ObjectPath)
}

/*
	CollectionDeleted (OUT ObjectPath collection);
*/

// SignalCollectionDeleted emits a signal that a collection was deleted
func (collection *Collection) SignalCollectionDeleted() {

	collection.Parent.Connection.Emit("/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service.CollectionDeleted",
		collection.ObjectPath)

	log.Infof("Emitted 'CollectionDeleted' signal for collection: %v", collection.ObjectPath)
}

/*
	CollectionChanged (OUT ObjectPath collection);
*/

// SignalCollectionDeleted emits a signal that a collection has changed
func (collection *Collection) SignalCollectionChanged() {

	collection.Parent.Connection.Emit("/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service.CollectionChanged",
		collection.ObjectPath)

	log.Infof("Emitted 'CollectionChanged' signal for collection: %v", collection.ObjectPath)
}

////////////////////////////// Properties //////////////////////////////

// GetProperty returns given dbus property value
func (collection *Collection) GetProperty(name string) (dbus.Variant, error) {

	busObject := collection.Parent.Connection.Object("org.freedesktop.secrets",
		collection.ObjectPath)

	variant, err := busObject.GetProperty("org.freedesktop.Secret.Collection." + name)

	if err != nil {
		return dbus.MakeVariant(nil), fmt.Errorf("error getting property '%s'. Error: %v", name, err)
	}

	collection.Properties[name] = variant

	return variant, nil
}

// SetProperty sets given dbus property name to given value
func (collection *Collection) SetProperty(name string, value interface{}) {

	collection.DbusProperties.SetMust("org.freedesktop.Secret.Collection",
		name, value)

	collection.DataMutex.Lock()
	collection.Properties[name] = dbus.MakeVariant(value)
	collection.DataMutex.Unlock()
}

// SetProperties processes raw properties and sets collection.Properties
func (collection *Collection) SetProperties(properties map[string]dbus.Variant) {

	const base string = "org.freedesktop.Secret.Collection."
	processedProperties := make(map[string]dbus.Variant)

	collection.rawProperties = properties // keep a copy of original properties

	for k, v := range properties {
		if len(k) < 35 || k[:34] != base {
			continue
		}
		key := strings.TrimSpace(k[34:])
		if len(key) == 0 ||
			key == "Items" ||
			key == "Locked" ||
			key == "Created" ||
			key == "Modified" {
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

	collection.Properties = processedProperties
	collection.Label = collection.Properties["Label"].Value().(string)
}

// UpdateModified updated 'Modified' dbus property of this collection
func (collection *Collection) UpdateModified() {

	collection.DataMutex.Lock()
	collection.Modified = uint64(time.Now().Unix())
	collection.DataMutex.Unlock()

	collection.DbusProperties.SetMust("org.freedesktop.Secret.Collection",
		"Modified", collection.Modified)
}

// UpdatePropertyCollections updates dbus property of this collection's items
func (collection *Collection) UpdatePropertyCollectionItems() {

	var items []string

	collection.ItemsMutex.RLock()
	for _, item := range collection.Items {
		items = append(items, string(item.ObjectPath))
	}
	collection.ItemsMutex.RUnlock()

	collection.DbusProperties.SetMust("org.freedesktop.Secret.Collection",
		"Items", items)
}

// Lock locks a collection and updates dbus 'Locked' and 'Modified' properties
func (collection *Collection) Lock() {
	collection.LockMutex.Lock()
	defer collection.LockMutex.Unlock()
	collection.Locked = true
	collection.SetProperty("Locked", true)
}

// Unlock unlocks a collection and updates dbus 'Locked' and 'Modified' properties
func (collection *Collection) Unlock() {
	collection.LockMutex.Lock()
	defer collection.LockMutex.Unlock()
	collection.Locked = false
	collection.SetProperty("Locked", false)
}
