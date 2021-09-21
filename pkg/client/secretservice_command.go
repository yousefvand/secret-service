package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

// OpenSession creates a session for encrypted or non-encrypted further communication
func (client *Client) SecretServiceCommand(
	command string, params string) (string, error) {

	var call *dbus.Call
	var err error

	call, err = client.Call("org.freedesktop.secrets", "/secretservice",
		"ir.remisa.SecretService", "Command", command, params)

	if err != nil {
		return "", errors.New("dbus call failed. Error: " + err.Error())
	}

	var result string
	err = call.Store(&result)

	if err != nil {
		return "", errors.New("type conversion failed in 'Command'. Error: " + err.Error())
	}

	return result, nil
}
