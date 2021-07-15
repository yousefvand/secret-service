package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

/*
	Lock ( IN Array<ObjectPath> objects,
	       OUT Array<ObjectPath> locked,
	       OUT ObjectPath Prompt);
*/

func TestClient_Lock(t *testing.T) {

	t.Run("Service Lock", func(t *testing.T) {

		ssClient, _ := client.New()
		session, _ := ssClient.OpenSession(client.Plain)

		collection1, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")
		collection2, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")
		collection3, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")

		// collection1 items
		secretApi11 := client.NewSecretApi()
		secretApi11.Session = session.ObjectPath
		item11, _, _ := collection1.CreateItem(map[string]dbus.Variant{"a": dbus.MakeVariant("b")}, secretApi11, true)

		secretApi12 := client.NewSecretApi()
		secretApi12.Session = session.ObjectPath
		item12, _, _ := collection1.CreateItem(map[string]dbus.Variant{"c": dbus.MakeVariant("d")}, secretApi12, true)

		// collection2 items
		secretApi21 := client.NewSecretApi()
		secretApi21.Session = session.ObjectPath
		item21, _, _ := collection2.CreateItem(map[string]dbus.Variant{"e": dbus.MakeVariant("f")}, secretApi21, true)

		secretApi22 := client.NewSecretApi()
		secretApi22.Session = session.ObjectPath
		item22, _, _ := collection2.CreateItem(map[string]dbus.Variant{"g": dbus.MakeVariant("h")}, secretApi22, true)

		// collection3 items
		secretApi31 := client.NewSecretApi()
		secretApi31.Session = session.ObjectPath
		item31, _, _ := collection3.CreateItem(map[string]dbus.Variant{"j": dbus.MakeVariant("k")}, secretApi31, true)

		secretApi32 := client.NewSecretApi()
		secretApi32.Session = session.ObjectPath
		item32, _, _ := collection3.CreateItem(map[string]dbus.Variant{"m": dbus.MakeVariant("n")}, secretApi32, true)

		lockCandidates := []dbus.ObjectPath{
			collection1.ObjectPath,
			item11.ObjectPath,
			item12.ObjectPath,

			// collection2.ObjectPath,
			// item21.ObjectPath,
			item22.ObjectPath,

			collection3.ObjectPath,
			item31.ObjectPath,
			// item32.ObjectPath,
		}

		if Service.GetCollectionByPath(collection1.ObjectPath).Locked {
			t.Errorf("collection1 is locked at service side: %v", collection1.ObjectPath)
		}

		locked, prompt, err := ssClient.Lock(lockCandidates)

		if err != nil {
			t.Errorf("Service lock failed. Error: %v", err)
		}

		if prompt != "/" {
			t.Errorf("Service lock returned unknown prompt: %v", prompt)
		}

		if !Service.GetCollectionByPath(collection1.ObjectPath).Locked {
			t.Errorf("collection1 is not locked at service side: %v", collection1.ObjectPath)
		}

		if !collection1.Locked {
			t.Errorf("collection1 is not locked at client side: %v", collection1.ObjectPath)
		}

		if !Service.GetItemByPath(item11.ObjectPath).Locked {
			t.Errorf("item11 is not locked at service side: %v", item11.ObjectPath)
		}

		if !item11.Locked {
			t.Errorf("item11 is not locked at client side: %v", item11.ObjectPath)
		}

		if !Service.GetItemByPath(item12.ObjectPath).Locked {
			t.Errorf("item12 is not locked at service side: %v", item12.ObjectPath)
		}

		if !item12.Locked {
			t.Errorf("item12 is not locked at client side: %v", item12.ObjectPath)
		}

		if Service.GetCollectionByPath(collection2.ObjectPath).Locked {
			t.Errorf("collection2 is locked at service side: %v", collection2.ObjectPath)
		}

		if collection2.Locked {
			t.Errorf("collection2 is locked at client side: %v", collection2.ObjectPath)
		}

		if Service.GetItemByPath(item21.ObjectPath).Locked {
			t.Errorf("item21 is not locked at service side: %v", item21.ObjectPath)
		}

		if item21.Locked {
			t.Errorf("item21 is not locked at client side: %v", item21.ObjectPath)
		}

		if !Service.GetItemByPath(item22.ObjectPath).Locked {
			t.Errorf("item22 is not locked at service side: %v", item22.ObjectPath)
		}

		if !item22.Locked {
			t.Errorf("item22 is not locked at client side: %v", item22.ObjectPath)
		}

		if !Service.GetCollectionByPath(collection3.ObjectPath).Locked {
			t.Errorf("collection3 is not locked at service side: %v", collection3.ObjectPath)
		}

		if !collection3.Locked {
			t.Errorf("collection3 is not locked at client side: %v", collection3.ObjectPath)
		}

		if !Service.GetItemByPath(item31.ObjectPath).Locked {
			t.Errorf("item31 is not locked at service side: %v", item31.ObjectPath)
		}

		if !item31.Locked {
			t.Errorf("item31 is not locked at client side: %v", item31.ObjectPath)
		}

		if Service.GetItemByPath(item32.ObjectPath).Locked {
			t.Errorf("item32 is not locked at service side: %v", item32.ObjectPath)
		}

		if item32.Locked {
			t.Errorf("item32 is not locked at client side: %v", item32.ObjectPath)
		}

		if len(locked) != 6 {
			t.Errorf("expected 6 locked objects, got %d", len(locked))
		}

		// Property Locked (collection 1 & 3 are locked)

		propertyLockedVariant, propertyErr := collection1.GetProperty("Locked")

		if propertyErr != nil {
			t.Errorf("Cannot read 'Locked' property. Error: %v", err)
		}

		propertyLocked, ok := propertyLockedVariant.Value().(bool)

		if !ok {
			t.Errorf("Expected property 'Locked' to be of type 'bool', got: '%T'",
				propertyLockedVariant.Value())
		}

		// BUG: Probably a godbus bug
		if !propertyLocked {
			t.Log("Expected collection1 to be locked but it is not (godbus bug?)")
		}

		propertyLockedVariant, propertyErr = collection3.GetProperty("Locked")

		if propertyErr != nil {
			t.Errorf("Cannot read 'Locked' property. Error: %v", err)
		}

		propertyLocked, ok = propertyLockedVariant.Value().(bool)

		if !ok {
			t.Errorf("Expected property 'Locked' to be of type 'bool', got: '%T'",
				propertyLockedVariant.Value())
		}

		if !propertyLocked {
			t.Error("Expected collection3 to be locked but it is not")
		}

	})
}
