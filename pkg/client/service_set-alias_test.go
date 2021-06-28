package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

/*
	SetAlias ( IN String name,
	           IN ObjectPath collection);
*/

func TestClient_SetAlias(t *testing.T) {

	t.Run("Service SetAlias", func(t *testing.T) {

		ssClient, _ := client.New()

		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "before")

		if Service.GetCollectionByAlias("before") == nil {
			t.Error("There is no collection with alias 'before'")
		}

		ssClient.SetAlias("after", collection.ObjectPath)

		if Service.GetCollectionByAlias("before") != nil {
			t.Error("There is still a collection with alias 'before'")
		}

		if Service.GetCollectionByAlias("after") == nil {
			t.Error("There is no collection with 'after' alias")
		}

		ssClient.SetAlias("/", collection.ObjectPath)

		if Service.GetCollectionByAlias("after") != nil {
			t.Error("There is still a collection with alias 'after'")
		}

		if Service.GetCollectionByPath(collection.ObjectPath).Alias != "" {
			t.Errorf("Collection '%v' alias is not empty", collection.ObjectPath)
		}

	})
}
