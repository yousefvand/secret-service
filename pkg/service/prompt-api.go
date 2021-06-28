// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"fmt"
	"os/exec"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

/*
	API implementation of:
	org.freedesktop.Secret.Prompt
*/

/*
	Current version of this service doesn't use prompting at all.
	In future versions when using prompting there are two possibilities:
	1. Just notify the user to unlock service (via cli app)
	2. Open cli app, type the command and prompt user for password then close terminal
	In latter case use 'xdotool' instead of 'wmctrl'
*/

/////////////////////////////////// Methods ///////////////////////////////////

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Prompt >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Prompt (IN String window-id);
*/

// perform the prompt. A prompt necessary to complete an operation
// windowId: Platform specific window handle to use for showing the prompt
func (prompt *Prompt) Prompt(windowId string) *dbus.Error {

	log.WithFields(log.Fields{
		"interface":   "org.freedesktop.Secret.Prompt",
		"method":      "Prompt",
		"prompt path": prompt.ObjectPath,
		"windowId":    windowId,
	}).Trace("Method called by client")

	if !CommandExists("wmctrl") {
		log.Error("Performing prompt needs 'wmctrl' command to be installed on system")
		return ApiErrorNotSupported()
	}

	cmd := exec.Command("wmctrl", "-ai", windowId)

	if err := cmd.Run(); err != nil {
		return DbusErrorCallFailed(fmt.Sprintf("'wmctrl -ai %s' failed. Error: %v", windowId, err.Error()))
	}

	// TODO: prompt object on dbus

	return nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Prompt <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Dismiss >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Dismiss (void);
*/

// Dismiss dismisses the prompt
func (prompt *Prompt) Dismiss() *dbus.Error {

	log.WithFields(log.Fields{
		"interface":   "org.freedesktop.Secret.Prompt",
		"method":      "Dismiss",
		"prompt path": prompt.ObjectPath,
	}).Trace("Method called by client")

	if !CommandExists("wmctrl") {
		log.Error("Performing prompt needs 'wmctrl' command to be installed on system")
		prompt.SignalPromptCompleted(false, dbus.MakeVariant(""))
		return ApiErrorNotSupported()
	}

	cmd := exec.Command("wmctrl", "-ci", prompt.WindowId)

	if err := cmd.Run(); err != nil {
		prompt.SignalPromptCompleted(false, dbus.MakeVariant(""))
		return DbusErrorCallFailed(fmt.Sprintf("'wmctrl -ci %s' failed. Error: %v", prompt.WindowId, err.Error()))
	}

	prompt.SignalPromptCompleted(true, dbus.MakeVariant(""))
	return nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Dismiss <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/////////////////////////////////// Signals ///////////////////////////////////

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Completed >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Completed ( OUT Boolean dismissed,
	OUT Variant result);
*/

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Completed <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
