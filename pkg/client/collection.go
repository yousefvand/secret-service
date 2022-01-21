package client

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
)

// NewCollection creates and initialize a new collection and returns it
func NewCollection(parent *Client) (*Collection, error) {
	collection := &Collection{}
	collection.Parent = parent
	collection.LockMutex = new(sync.Mutex)
	collection.ItemsMutex = new(sync.RWMutex)
	collection.Items = make(map[string]*Item)
	collection.SignalChan = make(chan *dbus.Signal)

	err := parent.Connection.AddMatchSignal(
		dbus.WithMatchObjectPath("/org/freedesktop/secrets"),
		dbus.WithMatchInterface("org.freedesktop.Secret.Collection"),
		dbus.WithMatchSender("org.freedesktop.secrets"),
	)

	if err != nil {
		return nil, errors.New("cannot watch for signals. Error: " + err.Error())
	}

	parent.Connection.Signal(collection.SignalChan)

	return collection, nil
}

// RemoveItem removes an item from the collection
func (collection *Collection) RemoveItem(itemPath dbus.ObjectPath) error {
	collection.ItemsMutex.Lock()
	defer collection.ItemsMutex.Unlock()
	if _, ok := collection.Items[string(itemPath)]; !ok {
		return errors.New("No such item to remove: " + string(itemPath))
	}
	delete(collection.Items, string(itemPath))
	return nil
}

// AddItem adds given item to the collection
func (collection *Collection) AddItem(item *Item) error {
	collection.ItemsMutex.Lock()
	defer collection.ItemsMutex.Unlock()
	if _, ok := collection.Items[string(item.ObjectPath)]; ok {
		return errors.New("Item already exist: " + string(item.ObjectPath))
	}
	collection.Items[string(item.ObjectPath)] = item
	return nil
}

// Lock, locks a collection
func (collection *Collection) Lock() {
	collection.LockMutex.Lock()
	collection.Locked = true
	collection.LockMutex.Unlock()
}

// Unlock, unlocks a collection
func (collection *Collection) Unlock() {
	collection.LockMutex.Lock()
	collection.Locked = false
	collection.LockMutex.Unlock()
}

// GetItemByPath returns an item based on its path, otherwise null
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

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Signals >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// org.freedesktop.Secret.Collection signals
type CollectionSignal uint8

const (
	ItemCreated CollectionSignal = iota
	ItemDeleted
	ItemChanged
)

// WatchSignal watches for desired signal within a time period
// If signal is received it returns true, otherwise false
func (collection *Collection) WatchSignal(signal CollectionSignal, timeout ...time.Duration) (bool, error) {

	var signalName string

	signalTimeout := time.Second // default timeout
	if len(timeout) > 0 {
		signalTimeout = timeout[0]
	}

	switch signal {
	case ItemCreated:
		signalName = "ItemCreated"
	case ItemDeleted:
		signalName = "ItemDeleted"
	case ItemChanged:
		signalName = "ItemChanged"
	}

	select {
	case signal := <-collection.SignalChan:
		if signal.Name == "org.freedesktop.Secret.Collection."+signalName {
			return true, nil
		} else {
			return false, fmt.Errorf("expected 'org.freedesktop.Secret.Collection.%s' signal got: %s", signalName, signal.Name)
		}
	case <-time.After(signalTimeout):
		return false, fmt.Errorf("receiving '%s' signal timed out", signalName)
	}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Signals <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Properties >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// GetProperty returns given dbus property value
func (collection *Collection) GetProperty(name string) (dbus.Variant, error) {

	busObject := collection.Parent.Connection.Object("org.freedesktop.secrets",
		collection.ObjectPath)

	variant, err := busObject.GetProperty("org.freedesktop.Secret.Collection." + name)

	if err != nil {
		return dbus.MakeVariant(nil),
			fmt.Errorf("error getting property '%s'. Error: %v", name, err)
	}

	collection.Properties[name] = variant

	return variant, nil
}

// SetProperty sets given dbus property name to given value
func (collection *Collection) SetProperty(name string, value interface{}) error {

	// collection.DbusProperties.SetMust("org.freedesktop.Secret.Collection",
	// 	name, value)
	busObject := collection.Parent.Connection.Object("org.freedesktop.secrets",
		collection.ObjectPath)

	err := busObject.SetProperty("org.freedesktop.Secret.Collection."+name, value)

	if err != nil {
		return fmt.Errorf("cannot set property '%s': %v", name, err)
	}

	collection.Properties[name] = dbus.MakeVariant(value)

	// Always update Modified
	modified, err := collection.GetProperty("Modified")

	if err != nil {
		return fmt.Errorf("failed to read 'Modified' property. Error: %v", err)
	}

	collection.Modified = modified.Value().(uint64)

	return nil
}

////////////////////////////// Property Wrappers //////////////////////////////

// PropertyGetItems returns 'Items' property of the collection
func (collection *Collection) PropertyGetItems() ([]string, error) {

	variant, err := collection.GetProperty("Items")

	if err != nil {
		return []string{}, fmt.Errorf("failed to read 'Items' property. Error: %v", err)
	}

	pathitems, ok := variant.Value().([]dbus.ObjectPath)

	if !ok {
		return []string{}, fmt.Errorf("expected 'Items' to be of type '[]dbus.ObjectPath', got: '%T'",
			variant.Value())
	}

	var items []string
	for k := range pathitems {
		items = append(items, string(pathitems[k]))
	}

	sort.Strings(items)

	var collectionItems []string
	for k := range collection.Items {
		collectionItems = append(collectionItems, k)
	}

	sort.Strings(collectionItems)

	// WON'T FIX: This is OK. Not all collection items on dbus created by a single client
	// if !reflect.DeepEqual(collectionItems, items) {
	// 	panic(fmt.Sprintf("Collection 'Items' property is out of sync. Object: %v, dbus: %v",
	// 		collectionItems, items))
	// }

	return items, nil
}

// PropertyGetLocked returns 'Locked' property of the collection
func (collection *Collection) PropertyGetLocked() (bool, error) {

	variant, err := collection.GetProperty("Locked")

	if err != nil {
		return false, fmt.Errorf("failed to read 'Locked' property. Error: %v", err)
	}

	locked, ok := variant.Value().(bool)

	if !ok {
		return false, fmt.Errorf("expected 'Locked' to be of type 'bool', got: '%T'",
			variant.Value())
	}

	if collection.Locked != locked {
		panic(fmt.Sprintf("collection 'Locked' property is out of sync. Object: %v, dbus: %v",
			collection.Locked, locked))
	}

	return locked, nil
}

// PropertyGetLabel returns 'Label' property of the collection
func (collection *Collection) PropertyGetLabel() (string, error) {

	variant, err := collection.GetProperty("Label")

	if err != nil {
		return "", fmt.Errorf("failed to read 'Label' property. Error: %v", err)
	}

	label, ok := variant.Value().(string)

	if !ok {
		return "", fmt.Errorf("expected 'Label' to be of type 'string', got: '%T'", variant.Value())
	}

	if collection.Label != label {
		panic(fmt.Sprintf("Collection 'Label' property is out of sync. Object: %v, dbus: %v",
			collection.Label, label))
	}

	return label, nil
}

// PropertySetLabel changes 'Label' property of the collection to the given value
func (collection *Collection) PropertySetLabel(label string) error {

	label = strings.TrimSpace(label)

	err := collection.SetProperty("Label", label)

	if err != nil {
		return fmt.Errorf("failed to write 'Label' property. Error: %v", err)
	}

	collection.Label = label

	return nil
}

// PropertyCreated returns 'Created' property of the collection
func (collection *Collection) PropertyCreated() (uint64, error) {

	variant, err := collection.GetProperty("Created")

	if err != nil {
		return 0, fmt.Errorf("failed to read 'Created' property. Error: %v", err)
	}

	created, ok := variant.Value().(uint64)

	if !ok {
		return 0, fmt.Errorf("expected 'Created' to be of type 'uint64', got: '%T'",
			variant.Value())
	}

	return created, nil
}

// PropertyModified returns 'Modified' property of the collection
func (collection *Collection) PropertyModified() (uint64, error) {

	variant, err := collection.GetProperty("Modified")

	if err != nil {
		return 0, fmt.Errorf("failed to read 'Modified' property. Error: %v", err)
	}

	modified, ok := variant.Value().(uint64)

	if !ok {
		return 0, fmt.Errorf("expected 'Modified' to be of type 'uint64', got: '%T'",
			variant.Value())
	}

	return modified, nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Properties <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
