package service

import (
	"log"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

var PropsService *prop.Properties

// dbusService creates SecretService objects and interfaces on dbus
func dbusService(service *Service) {

	////////////////////////////// Methods //////////////////////////////
	/*
		OpenSession ( IN String algorithm,
		              IN Variant input,
		              OUT Variant output,
		              OUT ObjectPath result);
	*/
	openSession := []introspect.Arg{
		{
			Name:      "algorithm",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "input",
			Type:      "v",
			Direction: "in",
		},
		{
			Name:      "output",
			Type:      "v",
			Direction: "out",
		},
		{
			Name:      "result",
			Type:      "o",
			Direction: "out",
		},
	}

	/*
		CreateCollection ( IN Dict<String,Variant> properties,
		IN String alias,
		OUT ObjectPath collection,
		OUT ObjectPath prompt);
	*/
	createCollection := []introspect.Arg{
		{
			Name:      "properties",
			Type:      "a{sv}",
			Direction: "in",
		},
		{
			Name:      "alias",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "collection",
			Type:      "o",
			Direction: "out",
		},
		{
			Name:      "prompt",
			Type:      "o",
			Direction: "out",
		},
	}

	/*
		SearchItems ( IN Dict<String,String> attributes,
		OUT Array<ObjectPath> unlocked,
		OUT Array<ObjectPath> locked);
	*/
	searchItems := []introspect.Arg{
		{
			Name:      "attributes",
			Type:      "a{ss}",
			Direction: "in",
		},
		{
			Name:      "unlocked",
			Type:      "ao",
			Direction: "out",
		},
		{
			Name:      "locked",
			Type:      "ao",
			Direction: "out",
		},
	}

	/*
		Unlock ( IN Array<ObjectPath> objects,
		OUT Array<ObjectPath> unlocked,
		OUT ObjectPath prompt);
	*/
	unlock := []introspect.Arg{
		{
			Name:      "objects",
			Type:      "ao",
			Direction: "in",
		},
		{
			Name:      "unlocked",
			Type:      "ao",
			Direction: "out",
		},
		{
			Name:      "prompt",
			Type:      "o",
			Direction: "out",
		},
	}

	/*
		Lock ( IN Array<ObjectPath> objects,
		OUT Array<ObjectPath> locked,
		OUT ObjectPath Prompt);
	*/
	lock := []introspect.Arg{
		{
			Name:      "objects",
			Type:      "ao",
			Direction: "in",
		},
		{
			Name:      "locked",
			Type:      "ao",
			Direction: "out",
		},
		{
			Name:      "prompt",
			Type:      "o",
			Direction: "out",
		},
	}

	/*
		GetSecrets ( IN Array<ObjectPath> items,
		IN ObjectPath session,
		OUT Dict<ObjectPath,Secret> secrets);
	*/
	getSecrets := []introspect.Arg{
		{
			Name:      "items",
			Type:      "ao",
			Direction: "in",
		},
		{
			Name:      "session",
			Type:      "o",
			Direction: "in",
		},
		{
			Name:      "secrets",
			Type:      "a{o(oayays)}",
			Direction: "out",
		},
	}

	/*
		ReadAlias ( IN String name,
		OUT ObjectPath collection);
	*/
	readAlias := []introspect.Arg{
		{
			Name:      "name",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "collection",
			Type:      "o",
			Direction: "out",
		},
	}

	/*
		SetAlias ( IN String name,
		IN ObjectPath collection);
	*/
	setAlias := []introspect.Arg{
		{
			Name:      "name",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "collection",
			Type:      "o",
			Direction: "in",
		},
	}

	////////////////////////////// Signals //////////////////////////////

	/*
		CollectionCreated (OUT ObjectPath collection);
	*/
	collectionCreated := []introspect.Arg{
		{
			Name: "collection",
			Type: "o",
		},
	}

	/*
		CollectionDeleted (OUT ObjectPath collection);
	*/
	collectionDeleted := []introspect.Arg{
		{
			Name: "collection",
			Type: "o",
		},
	}

	/*
		CollectionChanged (OUT ObjectPath collection);
	*/
	collectionChanged := []introspect.Arg{
		{
			Name: "collection",
			Type: "o",
		},
	}

	////////////////////////////// Properties //////////////////////////////

	/*
		READ Array<ObjectPath> Collections ;
	*/

	propsSpec := map[string]map[string]*prop.Prop{
		"org.freedesktop.Secret.Service": {
			"Collections": {
				Value:    []dbus.ObjectPath{"/org/freedesktop/secrets/aliases/default"},
				Writable: false,
				Emit:     prop.EmitTrue,
			},
		},
	}

	var err error
	PropsService, err = prop.Export(service.Connection, "/org/freedesktop/secrets", propsSpec)
	if err != nil {
		log.Panicf("export 'Service' propsSpec failed: %v", err)
	}

	/////////////////////////////////// dbus ///////////////////////////////////

	introService := &introspect.Node{
		Name: "/org/freedesktop/secrets",
		Children: []introspect.Node{
			{
				Name: "aliases",
			},
			{
				Name: "collection",
			},
			{
				Name: "session",
			},
		},
		Interfaces: []introspect.Interface{
			introspect.IntrospectData, prop.IntrospectData,
			{
				Name: "org.freedesktop.Secret.Service", // interface name
				Methods: []introspect.Method{
					{
						Name: "OpenSession",
						Args: openSession,
					},
					{
						Name: "CreateCollection",
						Args: createCollection,
					},
					{
						Name: "SearchItems",
						Args: searchItems,
					},
					{
						Name: "Unlock",
						Args: unlock,
					},
					{
						Name: "Lock",
						Args: lock,
					},
					{
						Name: "GetSecrets",
						Args: getSecrets,
					},
					{
						Name: "ReadAlias",
						Args: readAlias,
					},
					{
						Name: "SetAlias",
						Args: setAlias,
					},
				},
				Signals: []introspect.Signal{
					{
						Name: "CollectionCreated",
						Args: collectionCreated,
					},
					{
						Name: "CollectionDeleted",
						Args: collectionDeleted,
					},
					{
						Name: "CollectionChanged",
						Args: collectionChanged,
					},
				},
				Properties: PropsService.Introspection("org.freedesktop.Secret.Service"),
			},
		},
	}

	service.Connection.Export(service, "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service")

	service.Connection.Export(introspect.NewIntrospectable(introService),
		"/org/freedesktop/secrets", "org.freedesktop.DBus.Introspectable")

}
