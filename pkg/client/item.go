package client

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"
)

// NewCollection creates and initialize a new collection and returns it
func NewItem(parent *Collection) *Item {
	item := &Item{}
	item.Parent = parent
	item.Locked = false
	item.Secret = NewSecret(item)
	item.LockMutex = new(sync.Mutex)
	item.LookupAttributesMutex = new(sync.RWMutex)
	// item.Created = uint64(time.Now().Unix())
	// item.Modified = item.Created
	// TODO: Update
	return item
}

// Lock, locks the item
func (item *Item) Lock() {
	item.LockMutex.Lock()
	item.Locked = true
	item.LockMutex.Unlock()
}

// Unlock, unlocks the item
func (item *Item) Unlock() {
	item.LockMutex.Lock()
	item.Locked = false
	item.LockMutex.Unlock()
}

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Properties >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// GetProperty returns given dbus property value
func (item *Item) GetProperty(name string) (dbus.Variant, error) {

	busObject := item.Parent.Parent.Connection.Object("org.freedesktop.secrets",
		item.ObjectPath)

	variant, err := busObject.GetProperty("org.freedesktop.Secret.Item." + name)

	if err != nil {
		return dbus.MakeVariant(nil),
			fmt.Errorf("error getting property '%s'. Error: %v", name, err)
	}

	return variant, nil
}

// SetProperty sets given dbus property name to given value
func (item *Item) SetProperty(name string, value interface{}) error {

	// collection.DbusProperties.SetMust("org.freedesktop.Secret.Collection",
	// 	name, value)
	busObject := item.Parent.Parent.Connection.Object("org.freedesktop.secrets",
		item.ObjectPath)

	err := busObject.SetProperty("org.freedesktop.Secret.Item."+name, value)

	if err != nil {
		return fmt.Errorf("cannot set property '%s': %v", name, err)
	}

	// Don't use invalid values
	// switch name {
	// case "Attributes":
	// 	if val, ok := value.(map[string]string); ok {
	// 		item.LookupAttributesMutex.Lock()
	// 		item.LookupAttributes = val
	// 		item.LookupAttributesMutex.Unlock()
	// 	}
	// case "Label":
	// 	if val, ok := value.(string); ok {
	// 		item.Label = val
	// 	}
	// }

	// Always update Modified
	modified, err := item.GetProperty("Modified")

	if err != nil {
		return fmt.Errorf("failed to read 'Modified' property. Error: %v", err)
	}

	item.Modified = modified.Value().(uint64)

	return nil
}

////////////////////////////// Property Wrappers //////////////////////////////

// PropertyGetLocked returns 'Locked' property of the item
func (item *Item) PropertyGetLocked() (bool, error) {

	variant, err := item.GetProperty("Locked")

	if err != nil {
		return false, fmt.Errorf("failed to read 'Locked' property. Error: %v", err)
	}

	locked, ok := variant.Value().(bool)

	if !ok {
		return false, fmt.Errorf("expected 'Locked' to be of type 'bool', got: '%T'",
			variant.Value())
	}

	if item.Locked != locked {
		panic(fmt.Sprintf("Item 'Locked' property is out of sync. Object: %v, dbus: %v",
			item.Locked, locked))
	}

	return locked, nil
}

// PropertyGeAttributes returns 'Attributes' property of the item
func (item *Item) PropertyGetAttributes() (map[string]string, error) {

	variant, err := item.GetProperty("Attributes")

	if err != nil {
		return map[string]string{},
			fmt.Errorf("failed to read 'Attributes' property. Error: %v", err)
	}

	attributes, ok := variant.Value().(map[string]string)

	if !ok {
		return map[string]string{},
			fmt.Errorf("expected 'Attributes' to be of type 'map[string]string', got: '%T'",
				variant.Value())
	}

	if !reflect.DeepEqual(item.LookupAttributes, attributes) {
		panic(fmt.Sprintf("Item 'Attributes' property is out of sync. Object: %v, dbus: %v",
			item.LookupAttributes, attributes))
	}

	return attributes, nil
}

// PropertySetAttributes changes 'Attributes' property of the item to the given value
func (item *Item) PropertySetAttributes(attributes map[string]string) error {

	err := item.SetProperty("Attributes", attributes)

	if err != nil {
		return fmt.Errorf("failed to write 'Attributes' property. Error: %v", err)
	}

	// result := make(map[string]string)

	// for k, v := range attributes {
	// 	result[strings.TrimSpace(k)] = strings.TrimSpace(v)
	// }

	item.LookupAttributesMutex.Lock()
	item.LookupAttributes = attributes
	item.LookupAttributesMutex.Unlock()

	return nil
}

// PropertyGetLabel returns 'Label' property of the item
func (item *Item) PropertyGetLabel() (string, error) {

	variant, err := item.GetProperty("Label")

	if err != nil {
		return "", fmt.Errorf("failed to read 'Label' property. Error: %v", err)
	}

	label, ok := variant.Value().(string)

	if !ok {
		return "", fmt.Errorf("expected 'Label' to be of type 'string', got: '%T'",
			variant.Value())
	}

	if item.Label != label {
		panic(fmt.Sprintf("Item 'Label' property is out of sync. Object: %v, dbus: %v",
			item.Label, label))
	}

	return label, nil
}

// PropertySetLabel changes 'Label' property of the item to the given value
func (item *Item) PropertySetLabel(label string) error {

	label = strings.TrimSpace(label)

	err := item.SetProperty("Label", label)

	if err != nil {
		return fmt.Errorf("failed to write 'Label' property. Error: %v", err)
	}

	item.Label = label

	return nil
}

// PropertyCreated returns 'Created' property of the item
func (item *Item) PropertyCreated() (uint64, error) {

	variant, err := item.GetProperty("Created")

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

// PropertyModified returns 'Modified' property of the item
func (item *Item) PropertyModified() (uint64, error) {

	variant, err := item.GetProperty("Modified")

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
