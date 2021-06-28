package service

import (
	"github.com/godbus/dbus/v5/introspect"
)

// dbusAddSession adds session on dbus at: '/org/freedesktop/secrets/session/SESSION_NAME'
func dbusAddSession(service *Service, session *Session) {

	dbusUpdateSessions(service)

	introSession := &introspect.Node{
		Name: string(session.ObjectPath),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name: "org.freedesktop.Secret.Session",
				Methods: []introspect.Method{
					{
						Name: "Close", /* Close (void); */
					},
				},
			},
		},
	}

	service.Connection.Export(session, session.ObjectPath, "org.freedesktop.Secret.Session")

	service.Connection.Export(introspect.NewIntrospectable(introSession), session.ObjectPath,
		"org.freedesktop.DBus.Introspectable")
}

// dbusUpdateSessions updates all sessions on dbus
func dbusUpdateSessions(service *Service) {

	children := []introspect.Node{}

	service.SessionsMutex.RLock()
	for _, v := range service.Sessions {
		children = append(children, introspect.Node{Name: v.CreateMethodFromPath("")})
	}
	service.SessionsMutex.RUnlock()

	service.Connection.Export(introspect.NewIntrospectable(&introspect.Node{
		Name:     "/org/freedesktop/secrets/session",
		Children: children,
	}), "/org/freedesktop/secrets/session", "org.freedesktop.DBus.Introspectable")

}
