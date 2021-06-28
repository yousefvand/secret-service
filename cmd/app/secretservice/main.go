// Cli tool for controlling secretserviced daemon
package main

import (
	"io"

	log "github.com/sirupsen/logrus"
)

func main() {
	// TODO: Control secretserviced by signals and dbus interface
	log.SetOutput(io.Discard)
	log.Print("Not Implemented!")
}
