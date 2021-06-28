package internal_test

import (
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/internal"
)

func Test_App(t *testing.T) {

	t.Run("App load", func(t *testing.T) {

		app := internal.NewApp()
		app.Load()
		if app.Config == nil {
			t.Error("No config after app.Load()")
		}
	})

	t.Run("App notify", func(t *testing.T) {

		app := internal.NewApp()
		app.Config = &internal.Config{}
		app.Config.Icon = "view-private" // or "flag"

		title := "Notification test"
		body := "Secret service provides secure ways of storing credentials"
		duration := time.Millisecond * 5000

		connection, err := dbus.SessionBus()
		if err != nil {
			t.Errorf("Connecting to dbus failed. Cannot send notification. Error: %v", err)
		}
		app.Service.Connection = connection

		app.Notify(title, body, duration)
	})
}
