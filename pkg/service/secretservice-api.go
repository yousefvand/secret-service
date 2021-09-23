// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"path/filepath"
	"time"

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
	case "export database":
		dbFile := filepath.Join(service.Home, time.Now().Format("2006.01.02-15:04:05")+"-"+"db.json")
		store := service.EncryptDatabase
		service.EncryptDatabase = false
		service.EncryptDatabase = store
		Marshal(service, dbFile)
		return "ok", nil
	default:
		return "unknown", nil
	}

}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Command <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
