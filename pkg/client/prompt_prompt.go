package client

import (
	"errors"
)

/*
	Prompt (IN String window-id);
*/

// Prompt performs the prompt. A prompt necessary to complete an operation
// windowId: Platform specific window handle to use for showing the prompt
func (prompt *Prompt) Prompt(windowId string) error {

	_, err := prompt.Parent.Call("org.freedesktop.secrets", prompt.ObjectPath,
		"org.freedesktop.Secret.Prompt", "Prompt")

	if err != nil {
		return errors.New("dbus call failed. Error: " + err.Error())
	}

	return nil
}
