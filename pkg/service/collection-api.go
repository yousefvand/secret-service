// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

/*
	API implementation of:
	org.freedesktop.Secret.Collection
*/

/////////////////////////////////// Methods ///////////////////////////////////

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Delete >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Delete (OUT ObjectPath prompt);
*/

// Delete removes the collection
func (c *Collection) Delete() (dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":       "org.freedesktop.Secret.Collection",
		"method":          "Delete",
		"collection path": c.ObjectPath,
	}).Trace("Method called by client")

	if c.Alias == "default" {
		return dbus.ObjectPath("/"), DbusErrorCallFailed("Cannot delete default collection")
	}

	c.Parent.RemoveCollection(c)
	c.SignalCollectionDeleted()
	c.Parent.UpdatePropertyCollections()

	return dbus.ObjectPath("/"), nil

}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Delete <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> SearchItems >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	SearchItems ( IN Dict<String,String> attributes,
	              OUT Array<ObjectPath> results);
*/

// SearchItems Searches for items in this collection matching the lookup attributes
func (c *Collection) SearchItems(
	attributes map[string]string) ([]dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":       "org.freedesktop.Secret.Collection",
		"method":          "Delete",
		"collection path": c.ObjectPath,
		"attributes":      attributes,
	}).Trace("Method called by client")

	var items []dbus.ObjectPath

	for _, item := range c.Items {
		if IsMapSubsetSingleMatch(item.LookupAttributes, attributes, c.ItemsMutex) {
			items = append(items, item.ObjectPath)
		}
	}

	return items, nil

}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< SearchItems <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CreateItem >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	CreateItem ( IN Dict<String,Variant> properties,
	             IN Secret secret,
	             IN Boolean replace,
	             OUT ObjectPath item,
	             OUT ObjectPath prompt);
*/

// creates an item (secret + lookup attributes + label) in a collection
func (c *Collection) CreateItem(properties map[string]dbus.Variant,
	secretApi SecretApi, replace bool) (dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":       "org.freedesktop.Secret.Collection",
		"method":          "CreateItem",
		"collection path": c.ObjectPath,
		"properties":      properties,
		"secretApi":       secretApi,
		"replace":         replace,
	}).Trace("Method called by client")

	if len(properties) == 0 { // FIXME: Is this allowed by API? Return an error
		log.Warn("Client asked to create an item with empty 'properties' (no Label, no Attributes)")
	}
	item := NewItem(c)
	item.SetProperties(properties)

	item.Secret.SecretApi = &secretApi
	item.ObjectPath = dbus.ObjectPath(string(c.ObjectPath) + "/" + UUID())

	epoch := Epoch()
	err := c.AddItem(item, replace, true, false, epoch, epoch, false)

	if err != nil {
		return dbus.ObjectPath("/"), dbus.ObjectPath("/"), ApiErrorNoSession()
	}

	c.UpdateModified()

	log.WithFields(log.Fields{
		"Label":            item.Label,
		"LookupAttributes": item.LookupAttributes,
		"Plain Secret":     item.Secret.PlainSecret,
		"SecretApi":        item.Secret.SecretApi,
	}).Tracef("New Item added to collection: %s", c.ObjectPath)

	item.SignalItemCreated()
	c.UpdatePropertyCollectionItems()

	return item.ObjectPath, dbus.ObjectPath("/"), nil
}
