package service

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	log "github.com/sirupsen/logrus"
)

var PropsItem *prop.Properties

// introspectItemByPath is a generalized helper function for creating items
func introspectItemByPath(item *Item, locked bool, created uint64, modified uint64) *introspect.Interface {

	////////////////////////////// Methods //////////////////////////////
	/*
		Delete ( OUT ObjectPath Prompt);
	*/
	delete := []introspect.Arg{
		{
			Name:      "prompt",
			Type:      "o",
			Direction: "out",
		},
	}

	/*
		GetSecret ( IN ObjectPath session,
		OUT Secret secret);
	*/
	getSecret := []introspect.Arg{
		{
			Name:      "session",
			Type:      "o",
			Direction: "in",
		},
		{
			Name:      "secret",
			Type:      "(oayays)",
			Direction: "out",
		},
	}

	/*
		SetSecret ( IN Secret secret);
	*/
	setSecret := []introspect.Arg{
		{
			Name:      "secret",
			Type:      "(oayays)",
			Direction: "in",
		},
	}

	////////////////////////////// Properties //////////////////////////////

	// Item property specifications.
	// locked: is item locked
	// created: item creation time in epoch
	// modified: item modification time in epoch
	propsSpec := func(item *Item, locked bool, created uint64, modified uint64) map[string]map[string]*prop.Prop {

		props := make(map[string]*prop.Prop)

		// non-standard properties except Label
		/*
			READWRITE String Label ;
		*/
		for k, v := range item.Properties {
			props[k] = &prop.Prop{
				Value:    v.Value(),
				Writable: true,
				Emit:     prop.EmitTrue,
				Callback: func(p *prop.Change) *dbus.Error {
					if p.Name == "Label" {
						item.DataMutex.Lock()
						item.Label = p.Value.(string)
						item.DataMutex.Unlock()
					}
					item.Properties[p.Name] = dbus.MakeVariant(p.Value)
					log.Infof("Property '%v' of item '%v' changed to: %v",
						p.Name, item.ObjectPath, p.Value)

					go func(item *Item) {
						item.DataMutex.Lock()
						item.Modified = Epoch()
						item.DataMutex.Unlock()
						item.DbusProperties.SetMust("org.freedesktop.Secret.Item",
							"Modified", item.Modified)
						item.SignalItemChanged()
						item.SaveData()
					}(item)

					return nil
				},
			}
		}

		item.DataMutex.Lock()
		item.Created = created   // use provided time
		item.Modified = modified // use provided time
		item.DataMutex.Unlock()

		/*
			READ Array<ObjectPath> Items ;
		*/
		props["Attributes"] = &prop.Prop{
			Value:    map[string]string{},
			Writable: true,
			Emit:     prop.EmitTrue,
			Callback: func(p *prop.Change) *dbus.Error {
				if attributes, ok := p.Value.(map[string]string); ok {
					item.DataMutex.Lock()
					item.LookupAttributes = attributes
					item.DataMutex.Unlock()
					log.Infof("Property '%v' of item '%v' changed to: %v",
						p.Name, item.ObjectPath, p.Value)

					go func(item *Item) {
						item.DataMutex.Lock()
						item.Modified = Epoch()
						item.DataMutex.Unlock()
						item.DbusProperties.SetMust("org.freedesktop.Secret.Item",
							"Modified", item.Modified)
						item.SignalItemChanged()
						item.SaveData()
					}(item)
				}

				return nil
			},
		}

		/*
			READ Boolean Locked ;
		*/
		props["Locked"] = &prop.Prop{
			Value:    locked,
			Writable: false,
			Emit:     prop.EmitTrue,
		}

		/*
			READ UInt64 Created ;
		*/
		props["Created"] = &prop.Prop{
			Value:    created,
			Writable: false,
			Emit:     prop.EmitTrue,
		}

		/*
			READ UInt64 Modified ;
		*/
		props["Modified"] = &prop.Prop{
			Value:    modified,
			Writable: false,
			Emit:     prop.EmitTrue,
		}

		return map[string]map[string]*prop.Prop{"org.freedesktop.Secret.Item": props}
	}(item, locked, created, modified)

	var err error
	PropsItem, err = prop.Export(item.Parent.Parent.Connection, item.ObjectPath, propsSpec)
	if err != nil {
		log.Panicf("export 'item' propsSpec failed: %v", err)
	}
	item.DbusProperties = PropsItem

	////////////////////////////// remove me //////////////////////////////

	return &introspect.Interface{
		Name: "org.freedesktop.Secret.Item", // use provided object path
		Methods: []introspect.Method{
			{
				Name: "Delete",
				Args: delete,
			},
			{
				Name: "GetSecret",
				Args: getSecret,
			},
			{
				Name: "SetSecret",
				Args: setSecret,
			},
		},
		Properties: PropsItem.Introspection("org.freedesktop.Secret.Item"),
	}

}

// add item on dbus at: '/org/freedesktop/secrets/collection/COLLECTION_NAME/ITEM_NAME'
func dbusAddItem(collection *Collection, item *Item,
	locked bool, created uint64, modified uint64) {

	connection := collection.Parent.Connection

	dbusUpdateItems(collection, locked, created, modified)

	introItem := &introspect.Node{
		Name: string(item.ObjectPath),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			*introspectItemByPath(item, locked, created, modified),
		},
	}

	connection.Export(item, item.ObjectPath, "org.freedesktop.Secret.Item")

	connection.Export(introspect.NewIntrospectable(introItem), item.ObjectPath,
		"org.freedesktop.DBus.Introspectable")

}

// update dbus collections after add/remove a collection
func dbusUpdateItems(collection *Collection, locked bool, created uint64, modified uint64) {

	connection := collection.Parent.Connection
	children := []introspect.Node{}

	collection.ItemsMutex.RLock()
	for _, v := range collection.Items {
		children = append(children, introspect.Node{Name: v.CreateMethodFromPath("")})
	}
	collection.ItemsMutex.RUnlock()

	introCollection := &introspect.Node{
		Name: string(collection.ObjectPath),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			*introspectCollectionByPath(collection, locked, created, modified),
		},
		Children: children,
	}

	connection.Export(introspect.NewIntrospectable(introCollection),
		collection.ObjectPath, "org.freedesktop.DBus.Introspectable")

}
