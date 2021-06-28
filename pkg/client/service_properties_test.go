package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

// org.freedesktop.Secret.Service Properties test
func Test_Service_Properties(t *testing.T) {

	/*
		READ Array<ObjectPath> Collections ;
	*/
	t.Run("Service Property - Collections", func(t *testing.T) {

		ssClient, _ := client.New()

		// Add collection

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

		collections, err := ssClient.PropertyGetCollections()

		if err != nil {
			t.Errorf("Cannot read 'Collections' property. Error: %v", err)
		}

		if contains, err := client.SliceContains(collections, "/org/freedesktop/secrets/aliases/default"); err != nil {
			if !contains {
				t.Error("no default collection in 'Collections' property: '/org/freedesktop/secrets/aliases/default'")
			}
		}

		contains, err := client.SliceContains(collections, string(collection.ObjectPath))

		if err != nil {
			t.Errorf("'SliceContains' failed. Error: %v", err)
		}

		if !contains {
			t.Errorf("collection is not in 'Collections' property: %s", collection.ObjectPath)
		}

		// Remove collection

		collectionPath, err := collection.Delete()

		if err != nil {
			t.Errorf("Collection 'Delete' failed. Error: %v", err)
		}

		if ssClient.HasCollection(collectionPath) {
			t.Errorf("collection exists at client side: %s", collectionPath)
		}

		if Service.HasCollection(collectionPath) {
			t.Errorf("collection exists at service side: %s", collectionPath)
		}

		collections, err = ssClient.PropertyGetCollections()

		if err != nil {
			t.Errorf("Cannot read 'Collections' property. Error: %v", err)
		}

		if contains, err := client.SliceContains(collections, "/org/freedesktop/secrets/aliases/default"); err != nil {
			if !contains {
				t.Error("no default collection in 'Collections' property: '/org/freedesktop/secrets/aliases/default'")
			}
		}

		contains, err = client.SliceContains(collections, collectionPath)

		if err != nil {
			t.Errorf("'SliceContains' failed. Error: %v", err)
		}

		if contains {
			t.Errorf("collection is still in 'Collections' property: %s", collectionPath)
		}

	})

}
