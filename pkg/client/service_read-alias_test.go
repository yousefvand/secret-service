package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

/*
	ReadAlias ( IN String name,
	            OUT ObjectPath collection);
*/

func TestClient_ReadAlias(t *testing.T) {

	t.Run("Service ReadAlias", func(t *testing.T) {

		ssClient, _ := client.New()
		collection1, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "test")
		ssClient.CreateCollection(map[string]dbus.Variant{}, "")
		ssClient.CreateCollection(map[string]dbus.Variant{}, "")

		defaultCollectionPath, err := ssClient.ReadAlias("default")

		if err != nil {
			t.Errorf("Service ReadAlias failed. Error: %v", err)
		}

		if defaultCollectionPath != "/org/freedesktop/secrets/aliases/default" {
			t.Errorf("Service ReadAlias for 'default' collection is wrong. Expected: '/org/freedesktop/secrets/aliases/default', got %v", defaultCollectionPath)
		}

		testCollectionPath, _ := ssClient.ReadAlias("test")

		if collection1.ObjectPath != testCollectionPath {
			t.Errorf("Expected collection path with alias 'test' to be %v, got %v",
				collection1.ObjectPath, testCollectionPath)
		}

	})
}
