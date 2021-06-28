package service

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

// dbusInitialize creates default dbus objects i.e. '/org'...
func dbusInitialize(connection *dbus.Conn) {

	connection.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: "/",
		Children: []introspect.Node{
			{
				Name: "org",
			},
			{
				Name: "secretservice",
			},
		},
	}), "/", "org.freedesktop.DBus.Introspectable")

	// TODO: Implement dbus interface for service
	connection.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: "/secretservice",
	}), "/secretservice", "org.freedesktop.DBus.Introspectable")

	connection.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: "/org",
		Children: []introspect.Node{
			{
				Name: "freedesktop",
			},
		},
	}), "/org", "org.freedesktop.DBus.Introspectable")

	connection.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: "/org/freedesktop",
		Children: []introspect.Node{
			{
				Name: "secrets",
			},
		},
	}), "/org/freedesktop", "org.freedesktop.DBus.Introspectable")

}
