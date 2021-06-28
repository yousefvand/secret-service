package client

import (
	"errors"
)

/*
	Close (void);
*/

// Close closes a session
func (session *Session) Close() error {

	client := session.Parent
	_, err := client.Call("org.freedesktop.secrets", session.ObjectPath,
		"org.freedesktop.Secret.Session", "Close")

	if err != nil {
		return errors.New("dbus call failed. Error: " + err.Error())
	}

	if err = session.Remove(); err != nil {
		return errors.New("Session Close failed. Error: " + err.Error())
	}

	return nil
}
