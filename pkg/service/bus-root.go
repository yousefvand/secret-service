package service

import (
	"os"
	"time"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

// getRootName gets 'org.freedesktop.secrets' name exclusively on session dbus
func getRootName(connection *dbus.Conn) {

	const name = "org.freedesktop.secrets"

	// try 15 *2 = 30 seconds to acquire 'org.freedesktop.secrets'
	for retry := 0; retry < 15; retry++ {

		time.Sleep(time.Second * 2)
		replyRequestName, errRequestName := connection.RequestName(name,
			dbus.NameFlagDoNotQueue)

		if errRequestName != nil {
			log.Panicf("Cannot acquire name '%s' on dbus. Error: %v", name, errRequestName)
		}

		if replyRequestName == dbus.RequestNameReplyPrimaryOwner {
			return
		}

		log.Infof("Retry #%d failed to acquire 'org.freedesktop.secrets' name on dbus", retry+1)
	}
	log.Errorf("'%s' name is already taken on dbus! Exiting...", name)
	os.Exit(5)

}
