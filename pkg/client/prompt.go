package client

import (
	"errors"
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
)

// NewPrompt creates and initialize a new prompt
func NewPrompt(parent *Client) (*Prompt, error) {
	prompt := &Prompt{}
	prompt.Parent = parent
	prompt.SignalChan = make(chan *dbus.Signal)

	err := parent.Connection.AddMatchSignal(
		dbus.WithMatchObjectPath("/org/freedesktop/secrets"),
		dbus.WithMatchInterface("org.freedesktop.Secret.Prompt"),
		dbus.WithMatchSender("org.freedesktop.secrets"),
	)

	if err != nil {
		return nil, errors.New("cannot watch for signals. Error: " + err.Error())
	}

	parent.Connection.Signal(prompt.SignalChan)

	return prompt, nil
}

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Signal >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// org.freedesktop.Secret.Prompt signal
type PromptSignal uint8

const (
	Completed PromptSignal = iota
)

// WatchSignal watches for desired signal within a time period
// If signal is received it returns true, otherwise false
func (prompt *Prompt) WatchSignal(signal PromptSignal, timeout ...time.Duration) (bool, error) {

	signalTimeout := time.Second // default timeout
	if len(timeout) > 0 {
		signalTimeout = timeout[0]
	}

	select {
	case signal := <-prompt.SignalChan:
		if signal.Name == "org.freedesktop.Secret.Prompt.Completed" {
			return true, nil
		} else {
			return false, fmt.Errorf("expected 'org.freedesktop.Secret.Prompt.Completed' signal got: %s", signal.Name)
		}
	case <-time.After(signalTimeout):
		return false, fmt.Errorf("receiving 'Completed' signal timed out")
	}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Signal <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
