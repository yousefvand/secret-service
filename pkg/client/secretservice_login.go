package client

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

/*
	Login ( IN  String serialnumber,
	        IN  Array<Byte> passwordhash,
					IN  Array<Byte> passwordhash_iv,
					OUT Array<Byte> cookie,
					OUT Array<Byte> cookie_iv
					OUT String result);
*/

// Set password for first time or change a password of service.Returns cookie, cookie iv, result, error
func (client *Client) SecretServiceLogin(serialnumber string,
	passwordhash []byte, passwordhash_iv []byte) ([]byte, []byte, string, error) {

	var call *dbus.Call
	var err error

	call, err = client.Call("org.freedesktop.secrets", "/secretservice",
		"ir.remisa.SecretService", "Login", serialnumber, passwordhash, passwordhash_iv)

	if err != nil {
		return nil, nil, "failed", errors.New("dbus call failed. Error: " + err.Error())
	}

	var cookie, cookie_iv []byte
	var result string
	err = call.Store(&cookie, &cookie_iv, &result)

	if err != nil {
		return nil, nil, "failed", errors.New("type conversion failed in 'Login'. Error: " + err.Error())
	}

	return cookie, cookie_iv, result, nil

}
