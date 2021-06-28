package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

func Test_Collection_Signals(t *testing.T) {

	t.Run("Collection Signal - ItemCreated", func(t *testing.T) {

		ssClient, _ := client.New()

		// open session
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session. Error: %v", err)
		}

		// get default collection
		collection, prompt, err := ssClient.CreateCollection(map[string]dbus.Variant{}, "default")
		if err != nil {
			t.Errorf("cannot get default collection. Error: %v", err)
		}
		if prompt != "/" {
			t.Errorf("wrong prompt for getting default collection: %v", prompt)
		}

		if collection.ObjectPath != "/org/freedesktop/secrets/aliases/default" {
			t.Errorf("Expected defalt path at: '/org/freedesktop/secrets/aliases/default', got: %v", collection.ObjectPath)
		}

		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Skype for Desktop MSA/remisa"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"account":    "remisa",
				"service":    "Skype for Desktop MSA",
				"xdg:schema": "org.freedesktop.Secret.Generic",
			}),
		}

		iv, cipherData, err := client.AesCBCEncrypt([]byte("Victoria1"), session.SymmetricKey)

		if err != nil {
			t.Errorf("encryption1 error: %v", err)
		}

		secretApi := client.NewSecretApi()
		secretApi.ContentType = "text/plain"
		secretApi.Session = session.ObjectPath
		secretApi.Parameters = iv
		secretApi.Value = cipherData

		// Add item
		item, itemPrompt, itemErr := collection.CreateItem(properties, secretApi, true)

		if itemErr != nil {
			t.Errorf("CreateItem1 failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem1: %v", itemPrompt)
		}

		if item.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item1 path: %v", item.ObjectPath)
		}

		if len(item.ObjectPath) != 73 {
			t.Errorf("wrong item1 path length. Expected 73, got: %v", len(item.ObjectPath))
		}

		signalReceived, err := collection.WatchSignal(client.ItemCreated)

		if err != nil {
			t.Errorf("Failed to watch 'ItemCreated' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'ItemCreated' signal timed out")
		}

	})

	t.Run("Collection Signal - ItemDeleted", func(t *testing.T) {

		ssClient, _ := client.New()

		// open session
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session. Error: %v", err)
		}

		// get default collection
		collection, prompt, err := ssClient.CreateCollection(map[string]dbus.Variant{}, "default")
		if err != nil {
			t.Errorf("cannot get default collection. Error: %v", err)
		}
		if prompt != "/" {
			t.Errorf("wrong prompt for getting default collection: %v", prompt)
		}

		if collection.ObjectPath != "/org/freedesktop/secrets/aliases/default" {
			t.Errorf("Expected defalt path at: '/org/freedesktop/secrets/aliases/default', got: %v", collection.ObjectPath)
		}

		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Skype for Desktop MSA/remisa"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"account":    "remisa",
				"service":    "Skype for Desktop MSA",
				"xdg:schema": "org.freedesktop.Secret.Generic",
			}),
		}

		iv, cipherData, err := client.AesCBCEncrypt([]byte("Victoria1"), session.SymmetricKey)

		if err != nil {
			t.Errorf("encryption1 error: %v", err)
		}

		secretApi := client.NewSecretApi()
		secretApi.ContentType = "text/plain"
		secretApi.Session = session.ObjectPath
		secretApi.Parameters = iv
		secretApi.Value = cipherData

		// Add item
		item, itemPrompt, itemErr := collection.CreateItem(properties, secretApi, true)

		if itemErr != nil {
			t.Errorf("CreateItem1 failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem1: %v", itemPrompt)
		}

		if item.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item1 path: %v", item.ObjectPath)
		}

		if len(item.ObjectPath) != 73 {
			t.Errorf("wrong item1 path length. Expected 73, got: %v", len(item.ObjectPath))
		}

		signalReceived, err := collection.WatchSignal(client.ItemCreated)

		if err != nil {
			t.Errorf("Failed to watch 'ItemCreated' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'ItemCreated' signal timed out")
		}

		// Delete the item
		_, err = item.Delete()

		if err != nil {
			t.Errorf("Failed to delete item. Error: %v", err)
		}

		signalReceived, err = collection.WatchSignal(client.ItemDeleted)

		if err != nil {
			t.Errorf("Failed to watch 'ItemDeleted' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'ItemDeleted' signal timed out")
		}

	})

	t.Run("Collection Signal - ItemChanged", func(t *testing.T) {

		ssClient, _ := client.New()

		// open session
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session. Error: %v", err)
		}

		// get default collection
		collection, prompt, err := ssClient.CreateCollection(map[string]dbus.Variant{}, "default")
		if err != nil {
			t.Errorf("cannot get default collection. Error: %v", err)
		}
		if prompt != "/" {
			t.Errorf("wrong prompt for getting default collection: %v", prompt)
		}

		if collection.ObjectPath != "/org/freedesktop/secrets/aliases/default" {
			t.Errorf("Expected defalt path at: '/org/freedesktop/secrets/aliases/default', got: %v", collection.ObjectPath)
		}

		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Skype for Desktop MSA/remisa"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"account":    "remisa",
				"service":    "Skype for Desktop MSA",
				"xdg:schema": "org.freedesktop.Secret.Generic",
			}),
		}

		iv, cipherData, err := client.AesCBCEncrypt([]byte("Victoria1"), session.SymmetricKey)

		if err != nil {
			t.Errorf("encryption1 error: %v", err)
		}

		secretApi := client.NewSecretApi()
		secretApi.ContentType = "text/plain"
		secretApi.Session = session.ObjectPath
		secretApi.Parameters = iv
		secretApi.Value = cipherData

		// Add item
		item, itemPrompt, itemErr := collection.CreateItem(properties, secretApi, true)

		if itemErr != nil {
			t.Errorf("CreateItem1 failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem1: %v", itemPrompt)
		}

		if item.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item1 path: %v", item.ObjectPath)
		}

		if len(item.ObjectPath) != 73 {
			t.Errorf("wrong item1 path length. Expected 73, got: %v", len(item.ObjectPath))
		}

		signalReceived, err := collection.WatchSignal(client.ItemCreated)

		if err != nil {
			t.Errorf("Failed to watch 'ItemCreated' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'ItemCreated' signal timed out")
		}

		// Change the item

		iv2, cipherData2, err := client.AesCBCEncrypt([]byte("Victoria2"), session.SymmetricKey)

		if err != nil {
			t.Errorf("encryption2 error: %v", err)
		}

		secretApi2 := client.NewSecretApi()
		secretApi2.ContentType = "text/plain"
		secretApi2.Session = session.ObjectPath
		secretApi2.Parameters = iv2
		secretApi2.Value = cipherData2

		err = item.SetSecret(secretApi2)

		if err != nil {
			t.Errorf("Failed to set secret for item. Error: %v", err)
		}

		signalReceived, err = collection.WatchSignal(client.ItemChanged)

		if err != nil {
			t.Errorf("Failed to watch 'ItemChanged' signal. Error: %v", err)
		}

		if !signalReceived {
			t.Error("receiving 'ItemChanged' signal timed out")
		}

	})

}
