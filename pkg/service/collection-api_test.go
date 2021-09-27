package service_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
	"github.com/yousefvand/secret-service/pkg/crypto"
)

func Test_Delete(t *testing.T) {

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

func TestCollection_SearchItems(t *testing.T) {

	t.Run("Collection SearchItems", func(t *testing.T) {

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
			t.Errorf("CreateItem 1 failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem 1: %v", itemPrompt)
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
			t.Errorf("CreateItem 2 failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem 2: %v", itemPrompt)
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
			t.Errorf("CreateItem 3 failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem 3: %v", itemPrompt)
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

		////////////////////////////// search for items //////////////////////////////

		items, err := collection.SearchItems(map[string]string{"a": "b", "e": "f"})

		if err != nil {
			t.Errorf("SearchItems failed. Error: %v", err)
		}

		if !contains(items, item1.ObjectPath) {
			t.Errorf("Search result doesn't contain item1: %v", item1.ObjectPath)
		}
		if !contains(items, item2.ObjectPath) {
			t.Errorf("Search result doesn't contain item2: %v", item2.ObjectPath)
		}

		items2, err := collection.SearchItems(map[string]string{"y": "z"})

		if err != nil {
			t.Errorf("SearchItems failed. Error: %v", err)
		}

		if !contains(items2, item2.ObjectPath) {
			t.Errorf("Search result doesn't contain item2: %v", item2.ObjectPath)
		}
		if !contains(items2, item3.ObjectPath) {
			t.Errorf("Search result doesn't contain item3: %v", item3.ObjectPath)
		}

	})

}

func contains(slice []dbus.ObjectPath, element dbus.ObjectPath) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func Test_CreateItem(t *testing.T) {

	t.Run("Collection CreateItem", func(t *testing.T) {

		ssClient, _ := client.New()

		// open first session
		session1, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session1. Error: %v", err)
		}

		serviceDefaultCollection := Service.GetCollectionByAlias("default")

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

		properties1 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Skype for Desktop MSA/remisa"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"account":    "remisa",
				"service":    "Skype for Desktop MSA",
				"xdg:schema": "org.freedesktop.Secret.Generic",
			}),
		}

		iv1, cipherData1, err1 := crypto.AesCBCEncrypt([]byte("Victoria1"), session1.SymmetricKey)

		if err1 != nil {
			t.Errorf("encryption1 error: %v", err1)
		}

		secretApi1 := client.NewSecretApi()
		secretApi1.ContentType = "text/plain"
		secretApi1.Session = session1.ObjectPath
		secretApi1.Parameters = iv1
		secretApi1.Value = cipherData1

		// Add first item
		item1, itemPrompt, itemErr := collection.CreateItem(properties1, secretApi1, true)

		if itemErr != nil {
			t.Errorf("CreateItem1 failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem1: %v", itemPrompt)
		}

		if item1.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item1 path: %v", item1.ObjectPath)
		}

		if len(item1.ObjectPath) != 73 {
			t.Errorf("wrong item1 path length. Expected 73, got: %v", len(item1.ObjectPath))
		}

		serviceItem1 := serviceDefaultCollection.GetItemByPath(item1.ObjectPath)

		if serviceItem1 == nil {
			t.Errorf("No such item1 at service side: %s", item1.ObjectPath)
		} else {
			if serviceItem1.Secret.PlainSecret != "Victoria1" {
				t.Errorf("Expected plan secret to be 'Victoria1', got '%s'", serviceItem1.Secret.PlainSecret)
			}
			if serviceItem1.Label != properties1["org.freedesktop.Secret.Item.Label"].Value().(string) {
				t.Errorf("Item1 Label at service side: %s, expected: %s", serviceItem1.Label,
					properties1["org.freedesktop.Secret.Item.Label"].Value().(string))
			}
			if account := serviceItem1.GetLookupAttribute("account"); account != "remisa" {
				t.Errorf("Item1 Attribute 'account' at service side: %s, expected: %s", account,
					properties1["org.freedesktop.Secret.Item.Attributes"].Value().(map[string]string)["account"])
			}
			if service := serviceItem1.GetLookupAttribute("service"); service != "Skype for Desktop MSA" {
				t.Errorf("Item1 Attribute 'service' at service side: %s, expected: %s", service,
					properties1["org.freedesktop.Secret.Item.Attributes"].Value().(map[string]string)["service"])
			}
		}

		// open second session
		session2, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session2. Error: %v", err)
		}

		properties2 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Skype for Desktop remisa"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"account":    "RTJ",
				"service":    "Skype for Desktop MSA2",
				"xdg:schema": "org.freedesktop.Secret.Generic",
			}),
		}

		iv2, cipherData2, err2 := crypto.AesCBCEncrypt([]byte("Victoria2"), session2.SymmetricKey)

		if err != nil {
			t.Errorf("encryption2 error: %v", err2)
		}

		secretApi2 := client.NewSecretApi()
		secretApi2.ContentType = "text/plain"
		secretApi2.Session = session2.ObjectPath
		secretApi2.Parameters = iv2
		secretApi2.Value = cipherData2

		// Add second item
		item2, itemPrompt, itemErr := collection.CreateItem(properties2, secretApi2, true)

		if itemErr != nil {
			t.Errorf("CreateItem2 failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem2: %v", itemPrompt)
		}
		if item2.ObjectPath[:41] != "/org/freedesktop/secrets/aliases/default/" {
			t.Errorf("wrong item2 path: %v", item1.ObjectPath)
		}
		if len(item2.ObjectPath) != 73 {
			t.Errorf("wrong item2 path length. Expected 73, got: %v", len(item1.ObjectPath))
		}

		serviceItem2 := serviceDefaultCollection.GetItemByPath(item2.ObjectPath)

		if serviceItem2 == nil {
			t.Errorf("No such item2 at service side: %s", item2.ObjectPath)
		} else {
			if serviceItem2.Secret.PlainSecret != "Victoria2" {
				t.Errorf("Expected plan secret to be 'Victoria2', got '%s'", serviceItem2.Secret.PlainSecret)
			}
			if serviceItem2.Label != properties2["org.freedesktop.Secret.Item.Label"].Value().(string) {
				t.Errorf("Item2 Label at service side: %s, expected: %s", serviceItem2.Label,
					properties2["org.freedesktop.Secret.Item.Label"].Value().(string))
			}
			if account := serviceItem2.GetLookupAttribute("account"); account != "RTJ" {
				t.Errorf("Item2 Attribute 'account' at service side: %s, expected: %s", account,
					properties2["org.freedesktop.Secret.Item.Attributes"].Value().(map[string]string)["account"])
			}
			if service := serviceItem2.GetLookupAttribute("service"); service != "Skype for Desktop MSA2" {
				t.Errorf("Item2 Attribute 'service' at service side: %s, expected: %s", service,
					properties2["org.freedesktop.Secret.Item.Attributes"].Value().(map[string]string)["service"])
			}
		}

	})
}
