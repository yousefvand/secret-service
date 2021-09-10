package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	Command ( IN   String serialnumber,
		        IN   Array<Byte> cookie,
						IN   Array<Byte> cookie_iv,
						IN   Array<Byte> command,
						IN   Array<Byte> command_iv,
						IN   Array<Byte> params,
						IN   Array<Byte> params_iv,
						OUT  Array<Byte> result,
						OUT  Array<Byte> result_iv,);
*/

// OpenSession creates a session for encrypted or non-encrypted further communication
func (client *Client) SecretServiceCommand(
	serialnumber string, cookie []byte, cookie_iv []byte,
	command []byte, command_iv []byte, params []byte, params_iv []byte) ([]byte, []byte, error) {

	var call *dbus.Call
	var err error

	call, err = client.Call("org.freedesktop.secrets", "/secretservice",
		"ir.remisa.SecretService", "Command", serialnumber, cookie, cookie_iv,
		command, command_iv, params, params_iv)

	if err != nil {
		return nil, nil, errors.New("dbus call failed. Error: " + err.Error())
	}

	var result, result_iv []byte
	err = call.Store(&result, &result_iv)

	if err != nil {
		return nil, nil, errors.New("type conversion failed in 'Command'. Error: " + err.Error())
	}

	return result, result_iv, nil
}
