package service

import (
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

// dbusSecretService creates dbusSecretService objects and interfaces on dbus
func dbusSecretService(service *Service) {

	////////////////////////////// Simple Command (No security) //////////////////////////////

	/*
		Command ( IN   String command,
							IN   String params,
							OUT  String result);
	*/
	command := []introspect.Arg{
		{
			Name:      "command",
			Type:      "s",
			Direction: "in",
		},
		{
			Name:      "params",
			Type:      "s",
			Direction: "in",
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
						Name: "Command",
						Args: command,
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

	service.Connection.Export(service, "/secretservice",
		"ir.remisa.SecretService")

	service.Connection.Export(introspect.NewIntrospectable(introSecretService),
		"/secretservice", "org.freedesktop.DBus.Introspectable")

}
