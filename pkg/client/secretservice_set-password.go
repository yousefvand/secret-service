package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	SetPassword ( IN  String      serialnumber
		            IN  Array<Byte> oldpassword,
								IN  Array<Byte> oldpassword_iv,
								IN  Array<Byte> newpassword,
								IN  Array<Byte> newpassword_iv,
								IN  Array<Byte> oldSalt,
								IN  Array<Byte> oldSalt_iv,
								IN  Array<Byte> newSalt,
								IN  Array<Byte> newSalt_iv
								OUT String result);
*/

// Set password for first time or change a password of service
func (client *Client) SecretServiceSetPassword(serialnumber string,
	oldPassword []byte, oldPassword_iv []byte,
	newPassword []byte, newPassword_iv []byte,
	oldSalt []byte, oldSalt_iv []byte,
	newSalt []byte, newSalt_iv []byte) (string, error) {

	var call *dbus.Call
	var err error

	call, err = client.Call("org.freedesktop.secrets", "/secretservice",
		"ir.remisa.SecretService", "SetPassword", serialnumber,
		oldPassword, oldPassword_iv,
		newPassword, newPassword_iv,
		oldSalt, oldSalt_iv,
		newSalt, newSalt_iv,
	)

	if err != nil {
		return "failed", errors.New("dbus call failed. Error: " + err.Error())
	}

	var result string
	err = call.Store(&result)

	if err != nil {
		return "failed", errors.New("type conversion failed in 'SetPassword'. Error: " + err.Error())
	}

	return result, nil

}
