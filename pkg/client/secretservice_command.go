package client

/*
	Command ( IN String command,
	          OUT String result);
*/

// OpenSession creates a session for encrypted or non-encrypted further communication
func (client *Client) SecretServiceCommand(
	serialnumber string, command string, params []string) (string, error) {

	// client should already has obtained a CLI serialnumber
	if client.CliSession == nil || client.CliSession.SerialNumber == "" {
		panic("Client doesn't have a CLI serialnumber")
	}

	// use symmetric key to encrypt command and params and call dbus method

	return "", nil
}
