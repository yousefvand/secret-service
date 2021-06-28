// client package for 'secret service' as described at:
// http://standards.freedesktop.org/secret-service
package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

func TestNew(t *testing.T) {

	t.Run("client New", func(t *testing.T) {

		ssClient, err := client.New()
		if err != nil {
			t.Errorf("client cannot connect to session dbus. Error: %v", err)
			return
		}

		if ssClient.Connection == nil {
			t.Error("client connection is null")
		}

		if ssClient.Connected() == false {
			t.Error("client is not connected to session dbus")
		}

		/*

			err = client.Disconnect()
			if err != nil {
				t.Errorf("client cannot disconnect from session dbus. Error: %v", err.Error())
			}

			if client.Connected() == true {
				t.Error("client is connected to session dbus after Disconnect()")
			}

			connection, _ := dbus.SessionBus()
			client.Connection = connection

		*/

		// client.Lock()
		// defer client.Unlock()

	})

	t.Run("client call", func(t *testing.T) {

		ssClient, _ := client.New()
		call, err := ssClient.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
			"org.freedesktop.DBus.Introspectable", "Introspect")

		if err != nil {
			t.Errorf("dbus call failed. Error: %v", err.Error())
		}

		var output dbus.Variant
		err = call.Store(&output)

		if err != nil {
			t.Errorf("type conversion failed. Error: %v", err.Error())
		}

		xml := output.Value().(string)

		if len(xml) < 1 {
			t.Error("client method call failed")
		}

		_, err = ssClient.Call("org.freedesktop.secrets", "org/freedesktop/secrets",
			"org.freedesktop.DBus.Introspectable", "Introspect")

		if err == nil {
			t.Error("Invalid dbusPath didn't raise error")
		}

		_, err = ssClient.Call("org.freedesktop.secrets", "/org/freedesktop/secrets/",
			"org.freedesktop.DBus.Introspectable", "Introspect")

		if err == nil {
			t.Error("Invalid dbusPath didn't raise error")
		}

		_, err = ssClient.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
			".org.freedesktop.DBus.Introspectable", "Introspect")

		if err == nil {
			t.Error("Invalid dbus interface didn't raise error")
		}

		_, err = ssClient.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
			"org.freedesktop.DBus.Introspectable.", "Introspect")

		if err == nil {
			t.Error("Invalid dbus interface didn't raise error")
		}

		_, err = ssClient.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
			"org.freedesktop", "Introspect")

		if err == nil {
			t.Error("Invalid dbus interface didn't raise error")
		}

	})
}
