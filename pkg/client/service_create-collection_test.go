package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

/*
CreateCollection ( IN Dict<String,Variant> properties,
                   IN String alias,
                   OUT ObjectPath collection,
                   OUT ObjectPath prompt);
*/

func TestClient_CreateCollection(t *testing.T) {

	rawProperties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Collection.Label":  dbus.MakeVariant("MyCollection"),
		"org.freedesktop.Secret.Collection.Label1": dbus.MakeVariant("Test1"),
		"org.freedesktop.Secret.Collection.Label2": dbus.MakeVariant("Test2"),
	}

	t.Run("Service CreateCollection - empty properties", func(t *testing.T) {

		ssClient, _ := client.New()
		collection, promptPath, err := ssClient.CreateCollection(map[string]dbus.Variant{},
			"no-property")

		if err != nil {
			t.Errorf("CreateCollection failed. Error: %v", err)
		}

		if len(collection.ObjectPath) != 52 {
			t.Errorf("Invalid collection path length: %d", len(collection.ObjectPath))
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

	})

	t.Run("Service CreateCollection - default collection", func(t *testing.T) {

		alias := "default"
		ssClient, _ := client.New()

		collection, promptPath, err := ssClient.CreateCollection(rawProperties, alias)

		if err != nil {
			t.Errorf("CreateCollection for default collection failed. Error: %v", err)
		}

		if collection.ObjectPath != "/org/freedesktop/secrets/aliases/default" {
			t.Errorf("Invalid default collection path. Expected: /org/freedesktop/secrets/aliases/default, got: %s", collection.ObjectPath)
		}

		if promptPath != "/" {
			t.Errorf("Invalid prompt path for default collection: %s", promptPath)
		}

		if !ssClient.HasCollection(collection.ObjectPath) {
			t.Errorf("collection doesn't exist at client side: %s", collection.ObjectPath)
		}

		if !Service.HasCollection(collection.ObjectPath) {
			t.Errorf("collection doesn't exist at service side: %s", collection.ObjectPath)
		}

	})

	t.Run("Service CreateCollection - collection with alias", func(t *testing.T) {

		alias := "test"
		ssClient, _ := client.New()

		collection, promptPath, err := ssClient.CreateCollection(rawProperties, alias)

		if err != nil {
			t.Errorf("CreateCollection failed. Error: %v", err)
		}

		expectedPath := "/org/freedesktop/secrets/collection/" +
			rawProperties["org.freedesktop.Secret.Collection.Label"].Value().(string)

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

		// Get property Label
		label, err := collection.PropertyGetLabel()

		if err != nil {
			t.Errorf("Failed to read 'Label' property. Error: %v", err)
		}

		if label != "MyCollection" {
			t.Errorf("Wrong collection Label at client side. Expected 'MyCollection', got '%s'", label)
		}

		serviceSideLabel := Service.GetCollectionByPath(collection.ObjectPath).Label

		if serviceSideLabel != "MyCollection" {
			t.Errorf("Wrong collection Label at service side. Expected 'MyCollection', got '%s'", serviceSideLabel)
		}

		// Set property Label
		err = collection.PropertySetLabel("MyCollectionModified")

		if err != nil {
			t.Error(err)
		}

		label, err = collection.PropertyGetLabel()

		if err != nil {
			t.Errorf("Failed to read 'Label' property. Error: %v", err)
		}

		if label != "MyCollectionModified" {
			t.Errorf("Wrong collection Label. Expected 'MyCollection', got '%s'", label)
		}

	})

	t.Run("Service CreateCollection - collection with existing alias", func(t *testing.T) {

		alias := "something"
		ssClient, _ := client.New()

		collectionCountBefore := len(Service.Collections)
		collection, _, _ := ssClient.CreateCollection(rawProperties, alias)

		collection2, _, _ := ssClient.CreateCollection(rawProperties, alias)

		collectionCountAfter := len(Service.Collections)

		if collection.ObjectPath != collection2.ObjectPath {
			t.Errorf("More than one collection with the same alias '%s' at: %s",
				alias, collection.ObjectPath)
		}

		// This is not accurate, may need to be removed
		if collectionCountBefore != collectionCountAfter-1 {
			t.Errorf("Invalid collection count after creating with same alias: %s, %s",
				collection.ObjectPath, collection2.ObjectPath)
		}

	})

	t.Run("Service CreateCollection - collection with no alias", func(t *testing.T) {

		alias := ""
		ssClient, _ := client.New()

		collectionCountBefore := len(Service.Collections)
		collection, _, _ := ssClient.CreateCollection(rawProperties, alias)

		collection2, _, _ := ssClient.CreateCollection(rawProperties, alias)

		collectionCountAfter := len(Service.Collections)

		if collection.ObjectPath == collection2.ObjectPath {
			t.Errorf("Collections have same path: %s, %s",
				collection.ObjectPath, collection2.ObjectPath)
		}

		// This is not accurate, may need to be removed
		if collectionCountBefore != collectionCountAfter-2 {
			t.Errorf("Invalid collection count after creating with no alias: %s, %s",
				collection.ObjectPath, collection2.ObjectPath)
		}

	})

}
