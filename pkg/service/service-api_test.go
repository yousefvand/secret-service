package service_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
	"github.com/yousefvand/secret-service/pkg/crypto"
)

////////////////////////////// OpenSession //////////////////////////////

func Test_OpenSession(t *testing.T) {

	/*
		OpenSession ( IN String algorithm,
		IN Variant input,
		OUT Variant output,
		OUT ObjectPath result);
	*/

	t.Run("plain algorithm", func(t *testing.T) {

		ssClient, _ := client.New()
		session, err := ssClient.OpenSession(client.Plain)

		if err != nil {
			t.Errorf("OpenSession failed. Error: %v", err)
		}

		if session.ServicePublicKey != nil {
			t.Error("Unexpected public key for plain algorithm. Expected null.")
		}

		// i.e. "/org/freedesktop/secrets/session/uuid..." uuid is 32 character
		if session.ObjectPath[:33] != "/org/freedesktop/secrets/session/" {
			t.Errorf("invalid objectPath (should start with '/org/freedesktop/secrets/session/'): %s", session.ObjectPath)
		}

		if len(session.ObjectPath) != 65 {
			t.Errorf("invalid objectPath (length is not 32): %s", session.ObjectPath)
		}

		if !ssClient.HasSession(session.ObjectPath) {
			t.Errorf("session doesn't exist at client side: %s", session.ObjectPath)
		}

		if !Service.HasSession(session.ObjectPath) {
			t.Errorf("session doesn't exist at service side: %s", session.ObjectPath)
		}

	})

	t.Run("dh-ietf1024-sha256-aes128-cbc-pkcs7 algorithm", func(t *testing.T) {

		ssClient, _ := client.New()
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("OpenSession failed. Error: %v", err)
		}

		if len(session.ServicePublicKey) != 128 {
			t.Errorf("Unexpected public key length. Expected 128, got '%d'",
				len(session.ServicePublicKey))
		}

		// i.e. "/org/freedesktop/secrets/session/uuid..." uuid is 32 character
		if session.ObjectPath[:33] != "/org/freedesktop/secrets/session/" {
			t.Errorf("invalid objectPath (should start with '/org/freedesktop/secrets/session/'): %s", session.ObjectPath)
		}

		if len(session.ObjectPath) != 65 {
			t.Errorf("invalid objectPath (length is not 32): %s", session.ObjectPath)
		}

		if !ssClient.HasSession(session.ObjectPath) {
			t.Errorf("session doesn't exist at client side: %s", session.ObjectPath)
		}

		if !Service.HasSession(session.ObjectPath) {
			t.Errorf("session doesn't exist at service side: %s", session.ObjectPath)
		}

	})

	t.Run("unsupported algorithm", func(t *testing.T) {

		ssClient, _ := client.New()
		sessionCountBefore := len(Service.Sessions)
		_, err := ssClient.OpenSession(client.Unsupported)
		sessionCountAfter := len(Service.Sessions)

		if err == nil {
			t.Errorf("OpenSession didn't fail with unsupported algorithm. Error: %v", err)
		}

		if sessionCountAfter != sessionCountBefore {
			t.Errorf("Invalid session count after using unsupported algorithm. Expected : %d, got: %d",
				sessionCountBefore, sessionCountAfter)
		}

	})

}

////////////////////////////// CreateCollection //////////////////////////////

func Test_CreateCollection(t *testing.T) {

	properties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Collection.Label":  dbus.MakeVariant("MyCollection"),
		"org.freedesktop.Secret.Collection.Label1": dbus.MakeVariant("Test1"),
		"org.freedesktop.Secret.Collection.Label2": dbus.MakeVariant("Test2"),
	}

	t.Run("empty properties", func(t *testing.T) {

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

	t.Run("default collection", func(t *testing.T) {

		alias := "default"
		ssClient, _ := client.New()
		collection, promptPath, err := ssClient.CreateCollection(properties, alias)

		if err != nil {
			t.Errorf("CreateCollection failed. Error: %v", err)
		}

		if collection.ObjectPath != "/org/freedesktop/secrets/aliases/default" {
			t.Errorf("Invalid default collection path: %s", collection.ObjectPath)
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

	t.Run("collection with alias", func(t *testing.T) {

		alias := "test"
		ssClient, _ := client.New()

		collection, promptPath, err := ssClient.CreateCollection(properties, alias)

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

	})

	t.Run("collection with existing alias", func(t *testing.T) {

		alias := "something"
		ssClient, _ := client.New()
		collectionCountBefore := len(Service.Collections)
		collection, _, _ := ssClient.CreateCollection(properties, alias)
		collection2, _, _ := ssClient.CreateCollection(properties, alias)
		collectionCountAfter := len(Service.Collections)

		if collection.ObjectPath != collection2.ObjectPath {
			t.Errorf("More than one collection with the same alias '%s' at: %s",
				alias, collection.ObjectPath)
		}

		if collectionCountBefore != collectionCountAfter-1 {
			t.Errorf("Invalid collection count after creating with same alias: %s, %s",
				collection.ObjectPath, collection2.ObjectPath)
		}

	})

	t.Run("collection with no alias", func(t *testing.T) {

		alias := ""
		ssClient, _ := client.New()
		collectionCountBefore := len(Service.Collections)
		collection, _, _ := ssClient.CreateCollection(properties, alias)
		collection2, _, _ := ssClient.CreateCollection(properties, alias)
		collectionCountAfter := len(Service.Collections)

		if collection.ObjectPath == collection2.ObjectPath {
			t.Errorf("Collections have same path: %s, %s",
				collection.ObjectPath, collection2.ObjectPath)
		}

		if collectionCountBefore != collectionCountAfter-2 {
			t.Errorf("Invalid collection count after creating with no alias: %s, %s",
				collection.ObjectPath, collection2.ObjectPath)
		}

	})

}

func TestService_SearchItems(t *testing.T) {

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

func Test_Unlock(t *testing.T) {

	t.Run("Service Unlock", func(t *testing.T) {

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
			t.Errorf("collection2 is locked at service side: %v", collection1.ObjectPath)
		}

		if collection2.Locked {
			t.Errorf("collection2 is locked at client side: %v", collection1.ObjectPath)
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
			t.Errorf("collection3 is not locked at service side: %v", collection1.ObjectPath)
		}

		if !collection3.Locked {
			t.Errorf("collection3 is not locked at client side: %v", collection1.ObjectPath)
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

		unlocked, prompt2, err2 := ssClient.Unlock(lockCandidates)

		if err2 != nil {
			t.Errorf("Service unlock failed. Error: %v", err)
		}

		if prompt2 != "/" {
			t.Errorf("Service unlock returned unknown prompt: %v", prompt)
		}

		// all objects are unlocked

		// collection1
		if Service.GetCollectionByPath(collection1.ObjectPath).Locked {
			t.Errorf("collection1 is locked at service side: %v", collection1.ObjectPath)
		}

		if collection1.Locked {
			t.Errorf("collection1 is locked at client side: %v", collection1.ObjectPath)
		}

		if Service.GetItemByPath(item11.ObjectPath).Locked {
			t.Errorf("item11 is locked at service side: %v", item11.ObjectPath)
		}

		if item11.Locked {
			t.Errorf("item11 is locked at client side: %v", item11.ObjectPath)
		}

		if Service.GetItemByPath(item12.ObjectPath).Locked {
			t.Errorf("item12 is locked at service side: %v", item12.ObjectPath)
		}

		if item12.Locked {
			t.Errorf("item12 is locked at client side: %v", item12.ObjectPath)
		}

		// collection2
		if Service.GetCollectionByPath(collection2.ObjectPath).Locked {
			t.Errorf("collection2 is locked at service side: %v", collection1.ObjectPath)
		}

		if collection2.Locked {
			t.Errorf("collection2 is locked at client side: %v", collection1.ObjectPath)
		}

		if Service.GetItemByPath(item21.ObjectPath).Locked {
			t.Errorf("item21 is locked at service side: %v", item21.ObjectPath)
		}

		if item21.Locked {
			t.Errorf("item21 is locked at client side: %v", item21.ObjectPath)
		}

		if Service.GetItemByPath(item22.ObjectPath).Locked {
			t.Errorf("item22 is locked at service side: %v", item22.ObjectPath)
		}

		if item22.Locked {
			t.Errorf("item22 is locked at client side: %v", item22.ObjectPath)
		}

		// collection3
		if Service.GetCollectionByPath(collection3.ObjectPath).Locked {
			t.Errorf("collection3 is locked at service side: %v", collection1.ObjectPath)
		}

		if collection3.Locked {
			t.Errorf("collection3 is locked at client side: %v", collection1.ObjectPath)
		}

		if Service.GetItemByPath(item31.ObjectPath).Locked {
			t.Errorf("item31 is locked at service side: %v", item31.ObjectPath)
		}

		if item31.Locked {
			t.Errorf("item31 is locked at client side: %v", item31.ObjectPath)
		}

		if Service.GetItemByPath(item32.ObjectPath).Locked {
			t.Errorf("item32 is locked at service side: %v", item32.ObjectPath)
		}

		if item32.Locked {
			t.Errorf("item32 is locked at client side: %v", item32.ObjectPath)
		}

		if len(unlocked) != 6 {
			t.Errorf("expected 6 unlocked collections, got %d", len(unlocked))
		}

	})
}

func Test_GetSecrets(t *testing.T) {

	t.Run("Service GetSecrets", func(t *testing.T) {

		ssClient, _ := client.New()

		// open first session
		session1, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session1. Error: %v", err)
		}

		serviceDefaultCollection := Service.GetCollectionByAlias("default")

		// get client side default collection
		defaultCollection, prompt, err := ssClient.CreateCollection(map[string]dbus.Variant{}, "default")
		if err != nil {
			t.Errorf("cannot get default collection. Error: %v", err)
		}
		if prompt != "/" {
			t.Errorf("wrong prompt for getting default collection: %v", prompt)
		}

		if defaultCollection.ObjectPath != "/org/freedesktop/secrets/aliases/default" {
			t.Errorf("Expected defalt path at: '/org/freedesktop/secrets/aliases/default', got: %v", defaultCollection.ObjectPath)
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

		// Add first item (uses session1)
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

		// open second session
		session2, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session2. Error: %v", err)
		}

		properties2 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Skype for Desktop remisa"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"account":    "remisa2",
				"service":    "Skype for Desktop MSA",
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

		// Add second item (uses session2)
		item2, _, _ := defaultCollection.CreateItem(properties2, secretApi2, true)

		if len(item2.ObjectPath) != 73 {
			t.Errorf("wrong item2 path length. Expected 73, got: %v", len(item1.ObjectPath))
		}

		if serviceDefaultCollection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at service side: %s", item2.ObjectPath)
		}

		properties3 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Skype for Desktop/rem"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"account":    "remisa3",
				"service":    "Skype for Desktop",
				"xdg:schema": "org.freedesktop.Secret.Generic",
			}),
		}

		iv3, cipherData3, err3 := crypto.AesCBCEncrypt([]byte("Victoria3"), session2.SymmetricKey)

		if err3 != nil {
			t.Errorf("encryption3 error: %v", err3)
		}

		secretApi3 := client.NewSecretApi()
		secretApi3.ContentType = "text/plain"
		secretApi3.Session = session2.ObjectPath
		secretApi3.Parameters = iv3
		secretApi3.Value = cipherData3

		// Add third item (uses session2)
		item3, _, _ := defaultCollection.CreateItem(properties3, secretApi3, true)

		if len(item3.ObjectPath) != 73 {
			t.Errorf("wrong item3 path length. Expected 73, got: %v", len(item1.ObjectPath))
		}

		if serviceDefaultCollection.GetItemByPath(item3.ObjectPath) == nil {
			t.Errorf("No such item3 at service side: %s", item3.ObjectPath)
		}

		// getting secrets
		// open first session
		testSession, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session1. Error: %v", err)
		}

		secrets, err := ssClient.GetSecrets([]dbus.ObjectPath{
			item1.ObjectPath, item2.ObjectPath, item3.ObjectPath,
		},
			testSession.ObjectPath)

		if err != nil {
			t.Errorf("cannot 'GetSecrets'. Error: %v", err)
		}

		// item1
		if secretApi, ok := secrets[item1.ObjectPath]; ok {
			if secretApi.Session != testSession.ObjectPath {
				t.Errorf("session mismatch. Session: %s, SecretApi: %s",
					testSession.ObjectPath, secretApi.Session)
			}
			iv := secretApi.Parameters
			cipherData := secretApi.Value
			plainData, err := crypto.AesCBCDecrypt(iv, cipherData, testSession.SymmetricKey)
			if err != nil {
				t.Errorf("Error decrypting data. Error: %v", err)
			}
			if string(plainData) != "Victoria1" {
				t.Errorf("Expected: Victoria1, got %s", string(plainData))
			}
		} else {
			t.Errorf("item1 is not in the result: %v", item1.ObjectPath)
		}

		// item2
		if secretApi, ok := secrets[item2.ObjectPath]; ok {
			if secretApi.Session != testSession.ObjectPath {
				t.Errorf("session mismatch. Session: %s, SecretApi: %s",
					testSession.ObjectPath, secretApi.Session)
			}
			iv := secretApi.Parameters
			cipherData := secretApi.Value
			plainData, err := crypto.AesCBCDecrypt(iv, cipherData, testSession.SymmetricKey)
			if err != nil {
				t.Errorf("Error decrypting data. Error: %v", err)
			}
			if string(plainData) != "Victoria2" {
				t.Errorf("Expected: Victoria2, got %s", string(plainData))
			}
		} else {
			t.Errorf("item2 is not in the result: %v", item2.ObjectPath)
		}

		// item3
		if secretApi, ok := secrets[item3.ObjectPath]; ok {
			if secretApi.Session != testSession.ObjectPath {
				t.Errorf("session mismatch. Session: %s, SecretApi: %s",
					testSession.ObjectPath, secretApi.Session)
			}
			iv := secretApi.Parameters
			cipherData := secretApi.Value
			plainData, err := crypto.AesCBCDecrypt(iv, cipherData, testSession.SymmetricKey)
			if err != nil {
				t.Errorf("Error decrypting data. Error: %v", err)
			}
			if string(plainData) != "Victoria3" {
				t.Errorf("Expected: Victoria3, got %s", string(plainData))
			}
		} else {
			t.Errorf("item3 is not in the result: %v", item3.ObjectPath)
		}

	})
}

func Test_ReadAlias(t *testing.T) {

	ssClient, _ := client.New()
	collection1, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "test")
	ssClient.CreateCollection(map[string]dbus.Variant{}, "")
	ssClient.CreateCollection(map[string]dbus.Variant{}, "")

	defaultCollectionPath, err := ssClient.ReadAlias("default")

	if err != nil {
		t.Errorf("Service ReadAlias failed. Error: %v", err)
	}

	if defaultCollectionPath != "/org/freedesktop/secrets/aliases/default" {
		t.Errorf("Service ReadAlias for 'default' collection is wrong. Expected: '/org/freedesktop/secrets/aliases/default', got %v", defaultCollectionPath)
	}

	testCollectionPath, _ := ssClient.ReadAlias("test")

	if collection1.ObjectPath != testCollectionPath {
		t.Errorf("Expected collection path with alias 'test' to be %v, got %v",
			collection1.ObjectPath, testCollectionPath)
	}

}

func Test_SetAlias(t *testing.T) {

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
