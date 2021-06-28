package client

import "github.com/godbus/dbus/v5"

// NewSession creates and initialize a new session
func NewSession(parent *Client) *Session {
	session := &Session{}
	session.Parent = parent
	return session
}

// HasSession returns true if session exists otherwise false
func (client *Client) HasSession(sessionPath dbus.ObjectPath) bool {
	client.SessionsMutex.RLock()
	_, ok := client.Sessions[string(sessionPath)]
	client.SessionsMutex.RUnlock()
	return ok
}

// GetSessionByPath returns a session based on its path otherwise null
func (client *Client) GetSessionByPath(sessionPath dbus.ObjectPath) *Session {
	for _, session := range client.Sessions {
		if session.ObjectPath == sessionPath {
			return session
		}
	}
	return nil
}
