package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

func TestCollection_Delete(t *testing.T) {

	properties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Collection.Label":  dbus.MakeVariant("MyCollection"),
		"org.freedesktop.Secret.Collection.Label1": dbus.MakeVariant("Test1"),
		"org.freedesktop.Secret.Collection.Label2": dbus.MakeVariant("Test2"),
	}

	t.Run("Collection Delete", func(t *testing.T) {

		ssClient, _ := client.New()

		collection, promptPath, err := ssClient.CreateCollection(properties,
			"delete-test")

		if err != nil {
			t.Errorf("CreateCollection failed. Error: %v", err)
		}

		expectedPath := "/org/freedesktop/secrets/collection/" +
			properties["org.freedesktop.Secret.Collection.Label"].Value().(string)

		if string(collection.ObjectPath) != expectedPath {
			t.Errorf("Invalid collection path. Expected:  %s, got %s", expectedPath, collection.ObjectPath)
		}

		if promptPath != "/" {
			t.Errorf("Invalid prompt path: %s", promptPath)
		}

		if !ssClient.HasCollection(collection.ObjectPath) {
			t.Errorf("collection doesn't exist at client side: %s", collection.ObjectPath)
		}

		if !Service.HasCollection(collection.ObjectPath) {
			t.Errorf("collection doesn't exist at service side: %s", collection.ObjectPath)
		}

		// Delete collection
		collectionPath := collection.ObjectPath
		prompt, err := collection.Delete()

		if err != nil {
			t.Errorf("Cannot delete collection '%s'. Error: %v", collectionPath, err.Error())
		}

		if prompt != "/" {
			t.Errorf("Invalid prompt after deleting collection. Prompt: %s", prompt)
		}

		if ssClient.HasCollection(collectionPath) {
			t.Errorf("collection exists after deletion at client side: %s", collectionPath)
		}

		if Service.HasCollection(collectionPath) {
			t.Errorf("collection doesn't exist at service side: %s", collectionPath)
		}

	})

}
