package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

func Test_Service_Signals(t *testing.T) {

	t.Run("Service Signal - CollectionCreated", func(t *testing.T) {

		ssClient, _ := client.New()

		ssClient.CreateCollection(map[string]dbus.Variant{}, "")

		signalReceived, err := ssClient.WatchSignal(client.CollectionCreated)

		if err != nil {
			t.Errorf("Failed to watch 'CollectionCreated' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'CollectionCreated' signal timed out")
		}

		// get default collection
		ssClient.CreateCollection(map[string]dbus.Variant{}, "default")

		signalReceived, _ = ssClient.WatchSignal(client.CollectionCreated)

		if signalReceived {
			t.Error("Default collection cannot be created by client. No 'CollectionCreated' signal was expected")
		}

	})

	t.Run("Service Signal - CollectionDeleted", func(t *testing.T) {

		ssClient, _ := client.New()
		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")

		signalReceived, err := ssClient.WatchSignal(client.CollectionCreated)

		if err != nil {
			t.Errorf("Failed to watch 'CollectionCreated' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'CollectionCreated' signal timed out")
		}

		collection.Delete()

		signalReceived, err = ssClient.WatchSignal(client.CollectionDeleted)

		if err != nil {
			t.Errorf("Failed to watch 'CollectionDeleted' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'CollectionDeleted' signal timed out")
		}

	})

	t.Run("Service Signal - CollectionChanged", func(t *testing.T) {

		var rawProperties = map[string]dbus.Variant{
			"org.freedesktop.Secret.Collection.Label":  dbus.MakeVariant("RTJ"),
			"org.freedesktop.Secret.Collection.Color":  dbus.MakeVariant("White"),
			"org.freedesktop.Secret.Collection.Weight": dbus.MakeVariant(50),
		}

		alias := "before"
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

		signalReceived, err := ssClient.WatchSignal(client.CollectionCreated)

		if err != nil {
			t.Errorf("Failed to watch 'CollectionCreated' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'CollectionCreated' signal timed out")
		}

		// Set property Label
		err = collection.PropertySetLabel("Aamoo")

		if err != nil {
			t.Error(err)
		}

		label, err := collection.PropertyGetLabel()

		if err != nil {
			t.Errorf("Failed to read 'Label' property. Error: %v", err)
		}

		if label != "Aamoo" {
			t.Errorf("Wrong collection Label. Expected 'Aamoo', got '%s'", label)
		}

		signalReceived, err = ssClient.WatchSignal(client.CollectionChanged)
		if err != nil {
			t.Errorf("Failed to watch 'CollectionChanged' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'CollectionChanged' signal timed out")
		}

		// Change alias

		if Service.GetCollectionByAlias("before") == nil {
			t.Error("There is no collection with alias 'before'")
		}

		ssClient.SetAlias("after", collection.ObjectPath)

		signalReceived, err = ssClient.WatchSignal(client.CollectionChanged)
		if err != nil {
			t.Errorf("Failed to watch 'CollectionChanged' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'CollectionChanged' signal timed out")
		}

		if Service.GetCollectionByAlias("before") != nil {
			t.Error("There is still a collection with alias 'before'")
		}

		if Service.GetCollectionByAlias("after") == nil {
			t.Error("There is no collection with 'after' alias")
		}

		ssClient.SetAlias("/", collection.ObjectPath)

		signalReceived, err = ssClient.WatchSignal(client.CollectionChanged)
		if err != nil {
			t.Errorf("Failed to watch 'CollectionChanged' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'CollectionChanged' signal timed out")
		}

		if Service.GetCollectionByAlias("after") != nil {
			t.Error("There is still a collection with alias 'after'")
		}

		if Service.GetCollectionByPath(collection.ObjectPath).Alias != "" {
			t.Errorf("Collection '%v' alias is not empty", collection.ObjectPath)
		}

	})

}
