// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

/*
	API implementation of:
	org.freedesktop.Secret.Session
*/

/////////////////////////////////// Methods ///////////////////////////////////

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Close >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Close (void);
*/

// closes a session and removes its object from dbus
func (s *Session) Close() *dbus.Error {

	log.WithFields(log.Fields{
		"interface":    "org.freedesktop.Secret.Session",
		"method":       "Close",
		"session path": s.ObjectPath,
	}).Trace("Method called by client")

	s.Parent.RemoveSession(s)
	return nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Close <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
