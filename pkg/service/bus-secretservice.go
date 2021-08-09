package service

import (
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

// dbusSecretService creates dbusSecretService objects and interfaces on dbus
func dbusSecretService(secretservice *SecretService) {

	////////////////////////////// Methods //////////////////////////////
	/*
		OpenSession ( IN String algorithm,
		              IN Variant input,
		              OUT Variant output,
		              OUT String cookie);
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
			Type:      "s",
			Direction: "out",
		},
	}

	////////////////////////////// Signals //////////////////////////////

	/*
		ServiceLocked;
	*/
	serviceLocked := []introspect.Arg{
		{
			Name: "collection",
		},
	}

	/*
		ServiceUnlocked;
	*/
	serviceUnlocked := []introspect.Arg{
		{
			Name: "collection",
		},
	}

	/////////////////////////////////// dbus ///////////////////////////////////

	introSecretService := &introspect.Node{
		Name: "/secretservice",
		Interfaces: []introspect.Interface{
			introspect.IntrospectData, prop.IntrospectData,
			{
				Name: "ir.remisa.SecretService", // interface name
				Methods: []introspect.Method{
					{
						Name: "OpenSession",
						Args: openSession,
					},
				},
				Signals: []introspect.Signal{
					{
						Name: "ServiceLocked",
						Args: serviceLocked,
					},
					{
						Name: "serviceUnlocked",
						Args: serviceUnlocked,
					},
				},
			},
		},
	}

	secretservice.Parent.Connection.Export(secretservice, "/secretservice",
		"ir.remisa.SecretService")

	secretservice.Parent.Connection.Export(introspect.NewIntrospectable(introSecretService),
		"/secretservice", "org.freedesktop.DBus.Introspectable")

}
