package client_test

import (
	"reflect"
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

func TestNewCollection(t *testing.T) {

	t.Run("New collection", func(t *testing.T) {
		ssClient, err := client.New()
		if err != nil {
			t.Errorf("Cannot create new client. Error: %v", err)
		}
		collection, err := client.NewCollection(ssClient)

		if err != nil {
			t.Error(err)
		}

		if collection == nil {
			t.Error("collection is null")
		}

		if collection != nil && collection.Parent == nil {
			t.Error("collection parent is null")
		}

	})

	t.Run("HasCollection", func(t *testing.T) {

		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Collection.Label1": dbus.MakeVariant("Test1"),
			"org.freedesktop.Secret.Collection.Label2": dbus.MakeVariant("Test2"),
		}

		ssClient, _ := client.New()
		collection, _, _ := ssClient.CreateCollection(properties, "")

		if !ssClient.HasCollection(collection.ObjectPath) {
			t.Errorf("collection doesn't exist at client side: %s", collection.ObjectPath)
		}

		if !Service.HasCollection(collection.ObjectPath) {
			t.Errorf("collection doesn't exist at service side: %s", collection.ObjectPath)
		}

	})

	t.Run("GetCollectionByPath", func(t *testing.T) {

		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Collection.Label1": dbus.MakeVariant("Test1"),
			"org.freedesktop.Secret.Collection.Label2": dbus.MakeVariant("Test2"),
		}

		ssClient, _ := client.New()
		collection, _, _ := ssClient.CreateCollection(properties, "")

		if !reflect.DeepEqual(collection, ssClient.GetCollectionByPath(collection.ObjectPath)) {
			t.Errorf("collection doesn't match at client side: %s", collection.ObjectPath)
		}

		if ssClient.GetCollectionByPath("a/b/c") != nil {
			t.Error("Non existant collection exists!")
		}

	})

}
