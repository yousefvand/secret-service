package client_test

import (
	"reflect"
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

func TestNewSession(t *testing.T) {

	t.Run("New session", func(t *testing.T) {

		ssClient, err := client.New()
		if err != nil {
			t.Errorf("Cannot create new client. Error: %v", err)
		}
		session := client.NewSession(ssClient)

		if session == nil {
			t.Error("session is null")
		}

		if session != nil && session.Parent == nil {
			t.Error("session parent is null")
		}
	})

	t.Run("HasSession", func(t *testing.T) {

		ssClient, _ := client.New()
		session, _ := ssClient.OpenSession(client.Plain)

		if !ssClient.HasSession(session.ObjectPath) {
			t.Errorf("session doesn't exist at client side: %s", session.ObjectPath)
		}

		if !Service.HasSession(session.ObjectPath) {
			t.Errorf("session doesn't exist at service side: %s", session.ObjectPath)
		}

	})

	t.Run("GetSessionByPath", func(t *testing.T) {

		ssClient, _ := client.New()
		session, _ := ssClient.OpenSession(client.Plain)

		// FIXME
		if !reflect.DeepEqual(session, ssClient.GetSessionByPath(session.ObjectPath)) {
			t.Errorf("session doesn't match at client side: %s", session.ObjectPath)
		}

		if ssClient.GetSessionByPath("a/b/c") != nil {
			t.Error("Non existant session exists!")
		}

	})

}
