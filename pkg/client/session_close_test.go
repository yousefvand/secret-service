package client_test

import (
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

func TestSession_Close(t *testing.T) {

	t.Run("Close session", func(t *testing.T) {

		ssClient, _ := client.New()
		session, err := ssClient.OpenSession(client.Plain)

		if err != nil {
			t.Errorf("Failed to open session. Error: %v", err)
		}

		if ssClient.GetSessionByPath(session.ObjectPath) == nil {
			t.Errorf("Session doesn't exist at client side: %s", session.ObjectPath)
		}

		if Service.GetSessionByPath(session.ObjectPath) == nil {
			t.Errorf("Session doesn't exist at service side: %s", session.ObjectPath)
		}

		// close session
		err = session.Close()

		if err != nil {
			t.Errorf("Failed to close session. Error: %v", err)
		}

		if ssClient.GetSessionByPath(session.ObjectPath) != nil {
			t.Errorf("Session after 'Close' exist at client side: %s", session.ObjectPath)
		}

		if Service.GetSessionByPath(session.ObjectPath) != nil {
			t.Errorf("Session after 'Close' exist at service side: %s", session.ObjectPath)
		}

	})
}
