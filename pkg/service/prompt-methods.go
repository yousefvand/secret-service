// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

/*
	Current version of this service doesn't use prompting at all.
	In future versions when using prompting there are two possibilities:
	1. Just notify the user to unlock service (via cli app)
	2. Open cli app, type the command and prompt user for password then close terminal
	In latter case use 'xdotool' instead of 'wmctrl'
*/

// create and initialize a new collection
func NewPrompt(parent *Service) *Prompt {
	prompt := &Prompt{}
	prompt.Parent = parent
	// prompt.Update = parent.Update
	return prompt
}

func (prompt *Prompt) SignalPromptCompleted(dismissed bool, result dbus.Variant) {

	prompt.Parent.Connection.Emit("/org/freedesktop/secrets",
		"org.freedesktop.Secret.Prompt.Completed",
		dismissed, result)

	log.Infof("Emitted 'Completed' signal for prompt: %v, dismissed: %v, result: %v",
		prompt.ObjectPath, dismissed, result)
}
