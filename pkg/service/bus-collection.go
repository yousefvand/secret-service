// create, update dbus objects and interfaces
package service

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	log "github.com/sirupsen/logrus"
)

// Collection dbus properties
var PropsCollection *prop.Properties

// introspectCollectionByPath is a generalized helper function for creating collections
func introspectCollectionByPath(collection *Collection, locked bool, created uint64, modified uint64) *introspect.Interface {

	////////////////////////////// Methods //////////////////////////////
	/*
		Delete ( OUT ObjectPath prompt);
	*/
	delete := []introspect.Arg{
		{
			Name:      "prompt",
			Type:      "o",
			Direction: "out",
		},
	}

	/*
		SearchItems ( IN Dict<String,String> attributes,
		OUT Array<ObjectPath> results);
	*/
	searchItems := []introspect.Arg{
		{
			Name:      "attributes",
			Type:      "a{ss}",
			Direction: "in",
		},
		{
			Name:      "results",
			Type:      "ao",
			Direction: "out",
		},
	}

	/*
		CreateItem ( IN Dict<String,Variant> properties,
		IN Secret secret,
		IN Boolean replace,
		OUT ObjectPath item,
		OUT ObjectPath prompt);
	*/
	createItem := []introspect.Arg{
		{
			Name:      "properties",
			Type:      "a{sv}",
			Direction: "in",
		},
		{
			Name:      "secret",
			Type:      "(oayays)",
			Direction: "in",
		},
		{
			Name:      "replace",
			Type:      "b",
			Direction: "in",
		},
		{
			Name:      "item",
			Type:      "o",
			Direction: "out",
		},
		{
			Name:      "prompt",
			Type:      "o",
			Direction: "out",
		},
	}

	////////////////////////////// Signals //////////////////////////////

	/*
	   ItemCreated (OUT ObjectPath item);
	*/
	itemCreated := []introspect.Arg{
		{
			Name: "item",
			Type: "o",
		},
	}

	/*
	   ItemDeleted (OUT ObjectPath item);
	*/
	itemDeleted := []introspect.Arg{
		{
			Name: "item",
			Type: "o",
		},
	}

	/*
	   ItemChanged (OUT ObjectPath item);
	*/
	itemChanged := []introspect.Arg{
		{
			Name: "item",
			Type: "o",
		},
	}

	////////////////////////////// Properties //////////////////////////////

	// Collection property specifications.
	// locked: is collection locked
	// created: collection creation time in epoch
	// modified: collection modification time in epoch
	propsSpec := func(collection *Collection, locked bool, created uint64, modified uint64) map[string]map[string]*prop.Prop {

		props := make(map[string]*prop.Prop)

		// non-standard properties except Label
		/*
			READWRITE String Label ;
		*/
		for k, v := range collection.Properties {
			props[k] = &prop.Prop{
				Value:    v.Value(),
				Writable: true,
				Emit:     prop.EmitTrue,
				Callback: func(p *prop.Change) *dbus.Error {
					if p.Name == "Label" {
						collection.DataMutex.Lock()
						collection.Label = p.Value.(string)

						// FIXME: Change collection ObjectPath if available
						// if collection.Parent.GetCollectionByPath(
						// 	dbus.ObjectPath("/org/freedesktop/secrets/collection/"+collection.Label),
						// ) == nil {
						// 	collection.ObjectPath = dbus.ObjectPath(
						// 		"/org/freedesktop/secrets/collection/" + collection.Label)
						// }
						collection.DataMutex.Unlock()
					}
					collection.Properties[p.Name] = dbus.MakeVariant(p.Value)
					log.Infof("Property '%v' of collection '%v' changed to: %v",
						p.Name, collection.ObjectPath, p.Value)

					go func(collection *Collection) {
						collection.DataMutex.Lock()
						collection.Modified = Epoch()
						collection.DataMutex.Unlock()
						collection.DbusProperties.SetMust("org.freedesktop.Secret.Collection",
							"Modified", collection.Modified)
						collection.SignalCollectionChanged()
						collection.SaveData()
					}(collection)

					return nil
				},
			}
		}

		collection.DataMutex.Lock()
		collection.Created = created   // use provided time
		collection.Modified = modified // use provided time
		collection.DataMutex.Unlock()

		/*
			READ Array<ObjectPath> Items ;
		*/
		props["Items"] = &prop.Prop{
			Value:    []string{},
			Writable: false,
			Emit:     prop.EmitTrue,
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

		return map[string]map[string]*prop.Prop{"org.freedesktop.Secret.Collection": props}
	}(collection, locked, created, modified)

	var err error
	PropsCollection, err = prop.Export(collection.Parent.Connection, collection.ObjectPath, propsSpec)
	if err != nil {
		log.Panicf("export 'Collection' propsSpec failed: %v", err)
	}
	collection.DbusProperties = PropsCollection

	return &introspect.Interface{
		Name: "org.freedesktop.Secret.Collection", // use provided object path
		Methods: []introspect.Method{
			{
				Name: "Delete",
				Args: delete,
			},
			{
				Name: "SearchItems",
				Args: searchItems,
			},
			{
				Name: "CreateItem",
				Args: createItem,
			},
		},
		Signals: []introspect.Signal{
			{
				Name: "ItemCreated",
				Args: itemCreated,
			},
			{
				Name: "ItemDeleted",
				Args: itemDeleted,
			},
			{
				Name: "ItemChanged",
				Args: itemChanged,
			},
		},
		Properties: PropsCollection.Introspection("org.freedesktop.Secret.Collection"),
	}
}

// dbusDefaultCollection creates default collection at:
// '/org/freedesktop/secrets/aliases/default'
func dbusDefaultCollection(collection *Collection, locked bool, created uint64, modified uint64) {

	introCollection := &introspect.Node{
		Name: "/org/freedesktop/secrets/aliases/default",
		Interfaces: []introspect.Interface{
			introspect.IntrospectData, prop.IntrospectData,
			*introspectCollectionByPath(collection, locked, created, modified),
		},
	}

	collection.Parent.Connection.Export(collection, "/org/freedesktop/secrets/aliases/default",
		"org.freedesktop.Secret.Collection")

	collection.Parent.Connection.Export(introspect.NewIntrospectable(introCollection),
		"/org/freedesktop/secrets/aliases/default",
		"org.freedesktop.DBus.Introspectable")

}

// dbusAddCollection adds collection on dbus at:
// '/org/freedesktop/secrets/collection/COLLECTION_NAME'
func dbusAddCollection(collection *Collection, locked bool, created uint64, modified uint64) {

	dbusUpdateCollections(collection.Parent)

	introCollection := &introspect.Node{
		Name: string(collection.ObjectPath),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData, prop.IntrospectData,
			*introspectCollectionByPath(collection, locked, created, modified),
		},
	}

	collection.Parent.Connection.Export(collection, collection.ObjectPath,
		"org.freedesktop.Secret.Collection")

	collection.Parent.Connection.Export(introspect.NewIntrospectable(introCollection),
		collection.ObjectPath, "org.freedesktop.DBus.Introspectable")
}

// dbusUpdateCollections updates dbus collections after add/remove a collection
func dbusUpdateCollections(service *Service) {

	children := []introspect.Node{}

	service.CollectionsMutex.RLock()
	for _, v := range service.Collections {
		children = append(children, introspect.Node{Name: v.CreateMethodFromPath("")})
	}
	service.CollectionsMutex.RUnlock()

	service.Connection.Export(introspect.NewIntrospectable(&introspect.Node{
		Name:     "/org/freedesktop/secrets/collection",
		Children: children,
	}), "/org/freedesktop/secrets/collection", "org.freedesktop.DBus.Introspectable")

}
