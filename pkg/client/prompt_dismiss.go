package client

import (
	"errors"
)

/*
	Dismiss (void);
*/

// Dismiss dismisses the prompt
func (prompt *Prompt) Dismiss() error {

	_, err := prompt.Parent.Call("org.freedesktop.secrets", prompt.ObjectPath,
		"org.freedesktop.Secret.Prompt", "Dismiss")

	if err != nil {
		return errors.New("dbus call failed. Error: " + err.Error())
	}

	return nil
}
