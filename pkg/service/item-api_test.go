package service_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
	"github.com/yousefvand/secret-service/pkg/crypto"
)

func TestItem_Delete(t *testing.T) {

	t.Run("Item Delete", func(t *testing.T) {

		ssClient, _ := client.New()

		serviceDefaultCollection := Service.GetCollectionByAlias("default")

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

		////////////////////////////// item1 //////////////////////////////

		properties1 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("some item"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
			}),
		}

		iv1, cipherData1, err1 := crypto.AesCBCEncrypt([]byte("Victoria1"), session.SymmetricKey)

		if err1 != nil {
			t.Errorf("encryption1 error: %v", err1)
		}

		secretApi1 := client.NewSecretApi()
		secretApi1.ContentType = "text/plain"
		secretApi1.Session = session.ObjectPath
		secretApi1.Parameters = iv1
		secretApi1.Value = cipherData1

		// Add first item
		item1, itemPrompt, itemErr := collection.CreateItem(properties1, secretApi1, true)

		if itemErr != nil {
			t.Errorf("CreateItem failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", itemPrompt)
		}
		if item1.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item1 path: %v", item1.ObjectPath)
		}
		if len(item1.ObjectPath) != 73 {
			t.Errorf("wrong item1 path length. Expected 73, got: %v", len(item1.ObjectPath))
		}

		if collection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at client side: %s", item1.ObjectPath)
		}

		if serviceDefaultCollection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at service side: %s", item1.ObjectPath)
		}

		////////////////////////////// item2 //////////////////////////////

		properties2 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("another item"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"e": "f",
				"y": "z",
			}),
		}

		iv2, cipherData2, err2 := crypto.AesCBCEncrypt([]byte("Victoria2"), session.SymmetricKey)

		if err != nil {
			t.Errorf("encryption2 error: %v", err2)
		}

		secretApi2 := client.NewSecretApi()
		secretApi2.ContentType = "text/plain"
		secretApi2.Session = session.ObjectPath
		secretApi2.Parameters = iv2
		secretApi2.Value = cipherData2

		// Add second item
		item2, itemPrompt, itemErr := collection.CreateItem(properties2, secretApi2, true)

		if itemErr != nil {
			t.Errorf("CreateItem failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", itemPrompt)
		}
		if item2.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item2 path: %v", item2.ObjectPath)
		}
		if len(item2.ObjectPath) != 73 {
			t.Errorf("wrong item2 path length. Expected 73, got: %v", len(item2.ObjectPath))
		}

		if collection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at client side: %s", item2.ObjectPath)
		}

		if serviceDefaultCollection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at service side: %s", item2.ObjectPath)
		}

		////////////////////////////// item3 //////////////////////////////

		properties3 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("yet another item"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"j": "k",
				"o": "p",
				"y": "z",
			}),
		}

		iv3, cipherData3, err3 := crypto.AesCBCEncrypt([]byte("Victoria3"), session.SymmetricKey)

		if err != nil {
			t.Errorf("encryption3 error: %v", err3)
		}

		secretApi3 := client.NewSecretApi()
		secretApi3.ContentType = "text/plain"
		secretApi3.Session = session.ObjectPath
		secretApi3.Parameters = iv3
		secretApi3.Value = cipherData3

		// Add third item
		item3, itemPrompt, itemErr := collection.CreateItem(properties3, secretApi3, true)

		if itemErr != nil {
			t.Errorf("CreateItem failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", itemPrompt)
		}
		if item3.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item3 path: %v", item3.ObjectPath)
		}
		if len(item3.ObjectPath) != 73 {
			t.Errorf("wrong item3 path length. Expected 73, got: %v", len(item3.ObjectPath))
		}

		if collection.GetItemByPath(item3.ObjectPath) == nil {
			t.Errorf("No such item3 at client side: %s", item3.ObjectPath)
		}

		if serviceDefaultCollection.GetItemByPath(item3.ObjectPath) == nil {
			t.Errorf("No such item3 at service side: %s", item3.ObjectPath)
		}

		////////////////////////////// Delete items //////////////////////////////

		// item1
		item1Path, err := item1.Delete()

		if err != nil {
			t.Errorf("Item delete failed. Error: %v", err)
		}

		if collection.GetItemByPath(item1Path) != nil {
			t.Errorf("item1 still exists at client side: %s", item1Path)
		}

		if serviceDefaultCollection.GetItemByPath(item1Path) != nil {
			t.Errorf("item1 still exists at service side: %s", item1Path)
		}

		// item3
		item3Path, err := item3.Delete()

		if err != nil {
			t.Errorf("Item delete failed. Error: %v", err)
		}

		if collection.GetItemByPath(item3Path) != nil {
			t.Errorf("item3 still exists at client side: %s", item3Path)
		}

		if serviceDefaultCollection.GetItemByPath(item3Path) != nil {
			t.Errorf("item3 still exists at service side: %s", item3Path)
		}

		// we didn't delete item2 so it should exist
		if collection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at client side: %s", item2.ObjectPath)
		}

		if serviceDefaultCollection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at service side: %s", item2.ObjectPath)
		}

	})
}

func TestItem_GetSecret(t *testing.T) {

	t.Run("Item GetSecret", func(t *testing.T) {

		ssClient, _ := client.New()

		serviceDefaultCollection := Service.GetCollectionByAlias("default")

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

		////////////////////////////// item1 //////////////////////////////

		properties1 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("some item"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
			}),
		}

		iv1, cipherData1, err1 := crypto.AesCBCEncrypt([]byte("Victoria1"), session.SymmetricKey)

		if err1 != nil {
			t.Errorf("encryption1 error: %v", err1)
		}

		secretApi1 := client.NewSecretApi()
		secretApi1.ContentType = "text/plain"
		secretApi1.Session = session.ObjectPath
		secretApi1.Parameters = iv1
		secretApi1.Value = cipherData1

		// Add first item
		item1, itemPrompt, itemErr := collection.CreateItem(properties1, secretApi1, true)

		if itemErr != nil {
			t.Errorf("CreateItem failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", itemPrompt)
		}
		if item1.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item1 path: %v", item1.ObjectPath)
		}
		if len(item1.ObjectPath) != 73 {
			t.Errorf("wrong item1 path length. Expected 73, got: %v", len(item1.ObjectPath))
		}

		if collection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at client side: %s", item1.ObjectPath)
		}

		if serviceDefaultCollection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at service side: %s", item1.ObjectPath)
		}

		////////////////////////////// GetSecret //////////////////////////////

		secretApi, err := item1.GetSecret(session.ObjectPath)

		if err != nil {
			t.Errorf("GetSecret failed. Error: %v", err)
		}

		plainCipher, err := crypto.AesCBCDecrypt(secretApi.Parameters, secretApi.Value, session.SymmetricKey)

		if err != nil {
			t.Errorf("Decryption failed. Error: %v", err)
		}

		if string(plainCipher) != "Victoria1" {
			t.Errorf("Expected secret to be 'Victoria', got: %s", string(plainCipher))
		}

	})
}

func TestItem_SetSecret(t *testing.T) {

	t.Run("Item SetSecret", func(t *testing.T) {

		ssClient, _ := client.New()

		serviceDefaultCollection := Service.GetCollectionByAlias("default")

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

		////////////////////////////// item1 //////////////////////////////

		properties1 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("some item"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
			}),
		}

		iv1, cipherData1, err1 := crypto.AesCBCEncrypt([]byte("Victoria1"), session.SymmetricKey)

		if err1 != nil {
			t.Errorf("encryption1 error: %v", err1)
		}

		secretApi1 := client.NewSecretApi()
		secretApi1.ContentType = "text/plain"
		secretApi1.Session = session.ObjectPath
		secretApi1.Parameters = iv1
		secretApi1.Value = cipherData1

		// Add first item
		item1, itemPrompt, itemErr := collection.CreateItem(properties1, secretApi1, true)

		if itemErr != nil {
			t.Errorf("CreateItem failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", itemPrompt)
		}
		if item1.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item1 path: %v", item1.ObjectPath)
		}
		if len(item1.ObjectPath) != 73 {
			t.Errorf("wrong item1 path length. Expected 73, got: %v", len(item1.ObjectPath))
		}

		if collection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at client side: %s", item1.ObjectPath)
		}

		if serviceDefaultCollection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at service side: %s", item1.ObjectPath)
		}

		////////////////////////////// GetSecret //////////////////////////////

		originalSecretApi, err := item1.GetSecret(session.ObjectPath)

		if err != nil {
			t.Errorf("GetSecret failed. Error: %v", err)
		}

		originalPlainCipher, err := crypto.AesCBCDecrypt(originalSecretApi.Parameters,
			originalSecretApi.Value, session.SymmetricKey)

		if err != nil {
			t.Errorf("Decryption failed. Error: %v", err)
		}

		if string(originalPlainCipher) != "Victoria1" {
			t.Errorf("Expected secret to be 'Victoria1', got: %s", string(originalPlainCipher))
		}

		////////////////////////////// SetSecret (replace) //////////////////////////////

		iv2, cipherData2, err2 := crypto.AesCBCEncrypt([]byte("Victoria2"), session.SymmetricKey)

		if err2 != nil {
			t.Errorf("encryption2 error: %v", err2)
		}

		secretApi2 := client.NewSecretApi()
		secretApi2.ContentType = "text/plain"
		secretApi2.Session = session.ObjectPath
		secretApi2.Parameters = iv2
		secretApi2.Value = cipherData2

		err = item1.SetSecret(secretApi2)

		if err != nil {
			t.Errorf("SetSecret failed. Error: %v", err)
		}

		if item1.Secret.PlainSecret != "Victoria2" {
			t.Errorf("Expected secret to be 'Victoria2', got: %s", item1.Secret.PlainSecret)
		}

		////////////////////////////// Compare secrets //////////////////////////////

		changedSecretApi, err := item1.GetSecret(session.ObjectPath)

		if err != nil {
			t.Errorf("GetSecret failed. Error: %v", err)
		}

		changedPlainCipher, err := crypto.AesCBCDecrypt(changedSecretApi.Parameters,
			changedSecretApi.Value, session.SymmetricKey)

		if err != nil {
			t.Errorf("Decryption failed. Error: %v", err)
		}

		if string(changedPlainCipher) != "Victoria2" {
			t.Errorf("Expected secret to be 'Victoria2', got: %s", string(changedPlainCipher))
		}

	})
}
