package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
	"github.com/yousefvand/secret-service/pkg/crypto"
)

/*
	SearchItems ( IN Dict<String,String> attributes,
	              OUT Array<ObjectPath> unlocked,
	              OUT Array<ObjectPath> locked);
*/

func TestClient_SearchItems(t *testing.T) {

	t.Run("SearchItems", func(t *testing.T) {

		ssClient, _ := client.New()

		// get default collection
		defaultCollection, prompt, err := ssClient.CreateCollection(map[string]dbus.Variant{}, "default")
		if err != nil {
			t.Errorf("cannot get default collection. Error: %v", err)
		}
		if prompt != "/" {
			t.Errorf("wrong prompt for getting default collection: %v", prompt)
		}

		if defaultCollection.ObjectPath != "/org/freedesktop/secrets/aliases/default" {
			t.Errorf("Expected default path at: '/org/freedesktop/secrets/aliases/default', got: %v", defaultCollection.ObjectPath)
		}

		// Add items to default collection

		serviceDefaultCollection := Service.GetCollectionByAlias("default")

		// open session
		session1, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session1. Error: %v", err)
		}

		// Add first item to default collection
		properties1 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("SearchItems Test1"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"name":  "remisa",
				"age":   "26",
				"hobby": "math",
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

		item1, itemPrompt, itemErr := defaultCollection.CreateItem(properties1, secretApi1, true)

		if itemErr != nil {
			t.Errorf("CreateItem for item1 failed. Error: %v", itemErr)
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

		if serviceDefaultCollection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at service side: %s", item1.ObjectPath)
		}

		// Add second item to default collection
		properties2 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("SearchItems Test2"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"name":  "Remisa",
				"age":   "26",
				"hobby": "coding",
			}),
		}

		iv2, cipherData2, err2 := crypto.AesCBCEncrypt([]byte("Victoria2"), session1.SymmetricKey)

		if err2 != nil {
			t.Errorf("encryption2 error: %v", err2)
		}

		secretApi2 := client.NewSecretApi()
		secretApi2.ContentType = "text/plain"
		secretApi2.Session = session1.ObjectPath
		secretApi2.Parameters = iv2
		secretApi2.Value = cipherData2

		// Add second item to default collection
		item2, _, _ := defaultCollection.CreateItem(properties2, secretApi2, true)

		if len(item2.ObjectPath) != 73 {
			t.Errorf("wrong item2 path length. Expected 73, got: %v", len(item2.ObjectPath))
		}

		if serviceDefaultCollection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at service side: %s", item2.ObjectPath)
		}

		// Add a new collection
		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Collection.Label":  dbus.MakeVariant("SecretService"),
			"org.freedesktop.Secret.Collection.Type":   dbus.MakeVariant("Service"),
			"org.freedesktop.Secret.Collection.Target": dbus.MakeVariant("Linux"),
		}

		// non default collection
		collection, promptPath, err := ssClient.CreateCollection(properties, "Remisa")

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

		// Add first item to new collection
		properties3 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("SearchItems Test3"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"name":  "Ramiz",
				"age":   "6",
				"hobby": "Video Games",
			}),
		}

		iv3, cipherData3, err3 := crypto.AesCBCEncrypt([]byte("Victoria3"), session1.SymmetricKey)

		if err3 != nil {
			t.Errorf("encryption3 error: %v", err3)
		}

		secretApi3 := client.NewSecretApi()
		secretApi3.ContentType = "text/plain"
		secretApi3.Session = session1.ObjectPath
		secretApi3.Parameters = iv3
		secretApi3.Value = cipherData3

		item3, itemPrompt, itemErr := collection.CreateItem(properties3, secretApi3, true)

		if itemErr != nil {
			t.Errorf("CreateItem for item3 failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem 3: %v", itemPrompt)
		}
		if item3.ObjectPath[:50] != "/org/freedesktop/secrets/collection/SecretService/" {
			t.Errorf("wrong item3 path: %v", item3.ObjectPath)
		}
		if len(item3.ObjectPath) != 82 {
			t.Errorf("wrong item3 path length. Expected 82, got: %v", len(item3.ObjectPath))
		}

		if Service.GetCollectionByPath(collection.ObjectPath) == nil {
			t.Errorf("No such item3 at service side: %s", item3.ObjectPath)
		}

		// Add second item to new collection
		properties4 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("SearchItems Test4"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"name":  "Ramiz",
				"age":   "6",
				"hobby": "Playing with Aamoo",
			}),
		}

		iv4, cipherData4, err4 := crypto.AesCBCEncrypt([]byte("Victoria4"), session1.SymmetricKey)

		if err4 != nil {
			t.Errorf("encryption4 error: %v", err4)
		}

		secretApi4 := client.NewSecretApi()
		secretApi4.ContentType = "text/plain"
		secretApi4.Session = session1.ObjectPath
		secretApi4.Parameters = iv4
		secretApi4.Value = cipherData4

		item4, itemPrompt, itemErr := collection.CreateItem(properties4, secretApi4, true)

		if itemErr != nil {
			t.Errorf("CreateItem for item4 failed. Error: %v", itemErr)
		}
		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem 4: %v", itemPrompt)
		}
		if item4.ObjectPath[:50] != "/org/freedesktop/secrets/collection/SecretService/" {
			t.Errorf("wrong item4 path: %v", item4.ObjectPath)
		}
		if len(item3.ObjectPath) != 82 {
			t.Errorf("wrong item4 path length. Expected 82, got: %v", len(item4.ObjectPath))
		}

		if Service.GetCollectionByPath(collection.ObjectPath) == nil {
			t.Errorf("No such item4 at service side: %s", item4.ObjectPath)
		}

		////////////////////////////// SearchItems //////////////////////////////

		attributes := map[string]string{
			"name": "Ramiz",
			"age":  "6",
		}

		unlocked, locked, err := ssClient.SearchItems(attributes)

		if err != nil {
			t.Errorf("'SearchItems' failed. Error: %v", err)
		}

		if len(locked) > 0 {
			t.Errorf("Expected no locked items got: %v", locked)
		}

		if len(unlocked) != 2 {
			t.Errorf("Expected 2 results, got %v", len(unlocked))
		}

		attributes = map[string]string{
			"name":  "foo",
			"hobby": "math",
		}

		unlocked, locked, err = ssClient.SearchItems(attributes)

		if err != nil {
			t.Errorf("'SearchItems' failed. Error: %v", err)
		}

		if len(locked) > 0 {
			t.Errorf("Expected no locked items got: %v", locked)
		}

		if len(unlocked) != 0 {
			t.Errorf("Expected 0 results, got %v", len(unlocked))
		}

	})
}
