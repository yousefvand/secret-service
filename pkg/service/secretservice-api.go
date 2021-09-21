// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

/////////////////////////////////// Methods ///////////////////////////////////

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Command >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Command ( IN   String command,
						IN   String params,
						OUT  String result);
*/

// Command receives a command from CLI and runs it on daemon side
func (service *Service) Command(command string, params string) (string, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "ir.remisa.SecretService",
		"method":    "Command",
		"command":   command,
		"params":    params,
	}).Trace("Method called by client")

	switch command {
	case "ping":
		return "pong", nil
	case "foo":
		return "bar", nil
	default:
		return "unknown", nil
	}

}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Command <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
