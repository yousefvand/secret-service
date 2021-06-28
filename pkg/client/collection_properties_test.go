package client_test

import (
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

func Test_Collection_Properties(t *testing.T) {

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Items >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READ Array<ObjectPath> Items ;
	*/
	t.Run("Collection Property - Items", func(t *testing.T) {

		ssClient, _ := client.New()

		session, _ := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		// get default collection
		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "default")

		////////////////////////////// item1 //////////////////////////////

		properties1 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Label1"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
			}),
		}

		iv1, cipherData1, err1 := client.AesCBCEncrypt([]byte("Victoria1"), session.SymmetricKey)

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
			t.Errorf("CreateItem1 failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem1: %v", itemPrompt)
		}

		////////////////////////////// item2 //////////////////////////////

		properties2 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("Label2"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"w": "x",
				"y": "z",
			}),
		}

		iv2, cipherData2, err2 := client.AesCBCEncrypt([]byte("Victoria2"), session.SymmetricKey)

		if err2 != nil {
			t.Errorf("encryption2 error: %v", err2)
		}

		secretApi2 := client.NewSecretApi()
		secretApi2.ContentType = "text/plain"
		secretApi2.Session = session.ObjectPath
		secretApi2.Parameters = iv2
		secretApi2.Value = cipherData2

		// Add first item
		item2, itemPrompt, itemErr := collection.CreateItem(properties2, secretApi2, true)

		if itemErr != nil {
			t.Errorf("CreateItem2 failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem2: %v", itemPrompt)
		}

		////////////////////////////// check property //////////////////////////////

		items, err := collection.PropertyGetItems()

		if err != nil {
			t.Errorf("Cannot read 'Items' property of default collection. Error: %v", err)
		}

		// Check if item1 is there
		contains, err := client.SliceContains(items, string(item1.ObjectPath))

		if err != nil {
			t.Errorf("'SliceContains' failed. Error: %v", err)
		}

		if !contains {
			t.Errorf("item1 is not in 'Items' property: %s", item1.ObjectPath)
		}

		// Check if item2 is there
		contains, err = client.SliceContains(items, string(item2.ObjectPath))

		if err != nil {
			t.Errorf("'SliceContains' failed. Error: %v", err)
		}

		if !contains {
			t.Errorf("item1 is not in 'Items' property: %s", item2.ObjectPath)
		}

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Items <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Label >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READWRITE String Label ;
	*/
	t.Run("Collection Property - Label", func(t *testing.T) {

		ssClient, _ := client.New()

		InitialLabel := "LabelTest"

		rawProperties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Collection.Label":   dbus.MakeVariant(InitialLabel),
			"org.freedesktop.Secret.Collection.Custom1": dbus.MakeVariant("Test1"),
			"org.freedesktop.Secret.Collection.Custom2": dbus.MakeVariant("Test2"),
		}

		collection, _, _ := ssClient.CreateCollection(rawProperties, "")

		if string(collection.ObjectPath) != "/org/freedesktop/secrets/collection/"+InitialLabel {
			t.Errorf("Expected collection path to be: '%s', got: '%s'",
				"/org/freedesktop/secrets/collection/"+InitialLabel, string(collection.ObjectPath))
		}

		if collection.Label != InitialLabel {
			t.Errorf("Expected collection label to be '%s', got: '%s'",
				InitialLabel, collection.Label)
		}

		serviceCollection := Service.GetCollectionByPath(collection.ObjectPath)

		if serviceCollection.Label != InitialLabel {
			t.Errorf("Expected collection label to be '%s', got: '%s'",
				InitialLabel, serviceCollection.Label)
		}

		propertyLabel, err := collection.PropertyGetLabel()

		if err != nil {
			t.Error(err)
		}

		if propertyLabel != InitialLabel {
			t.Errorf("Expected collection label to be '%s', got: '%s'",
				InitialLabel, propertyLabel)
		}

		// Set Label
		secondaryLabel := "LabelTest2"
		err = collection.PropertySetLabel(secondaryLabel)

		if err != nil {
			t.Error(err)
		}

		// FIXME: Change collection path if available else use UUID

		if collection.Label != secondaryLabel {
			t.Errorf("Expected collection label to be '%s', got: '%s'",
				secondaryLabel, collection.Label)
		}

		serviceCollection = Service.GetCollectionByPath(collection.ObjectPath)

		if serviceCollection.Label != secondaryLabel {
			t.Errorf("Expected collection label to be '%s', got: '%s'",
				secondaryLabel, serviceCollection.Label)
		}

		propertyLabel, err = collection.PropertyGetLabel()

		if err != nil {
			t.Error(err)
		}

		if propertyLabel != secondaryLabel {
			t.Errorf("Expected collection label to be '%s', got: '%s'",
				secondaryLabel, propertyLabel)
		}

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Label <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Locked >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READ Boolean Locked ;
	*/
	t.Run("Collection Property - Locked", func(t *testing.T) {

		ssClient, _ := client.New()

		collection1, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")
		collection2, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")
		collection3, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")

		if collection1.Locked {
			t.Errorf("collection1 is initially locked at client side: %s", collection1.ObjectPath)
		}

		if Service.GetCollectionByPath(collection1.ObjectPath).Locked {
			t.Errorf("collection1 is initially locked at service side: %s", collection1.ObjectPath)
		}

		if collection2.Locked {
			t.Errorf("collection2 is initially locked at client side: %s", collection2.ObjectPath)
		}

		if Service.GetCollectionByPath(collection2.ObjectPath).Locked {
			t.Errorf("collection2 is initially locked at service side: %s", collection2.ObjectPath)
		}

		if collection3.Locked {
			t.Errorf("collection3 is initially locked at client side: %s", collection3.ObjectPath)
		}

		if Service.GetCollectionByPath(collection3.ObjectPath).Locked {
			t.Errorf("collection3 is initially locked at service side: %s", collection3.ObjectPath)
		}

		// lock collection1 and collection3
		lockCandidates := []dbus.ObjectPath{
			collection1.ObjectPath,
			// collection2.ObjectPath,
			collection3.ObjectPath,
		}

		locked, prompt, err := ssClient.Lock(lockCandidates)

		if err != nil {
			t.Errorf("Service lock failed. Error: %v", err)
		}

		if prompt != "/" {
			t.Errorf("Service lock returned unknown prompt: %s", prompt)
		}

		if !Service.GetCollectionByPath(collection1.ObjectPath).Locked {
			t.Errorf("collection1 is not locked at service side: %s", collection1.ObjectPath)
		}

		if !collection1.Locked {
			t.Errorf("collection1 is not locked at client side: %s", collection1.ObjectPath)
		}

		if Service.GetCollectionByPath(collection2.ObjectPath).Locked {
			t.Errorf("collection2 is locked at service side: %s", collection2.ObjectPath)
		}

		if collection2.Locked {
			t.Errorf("collection2 is locked at client side: %s", collection2.ObjectPath)
		}

		if !Service.GetCollectionByPath(collection3.ObjectPath).Locked {
			t.Errorf("collection3 is not locked at service side: %s", collection3.ObjectPath)
		}

		if !collection3.Locked {
			t.Errorf("collection3 is not locked at client side: %s", collection3.ObjectPath)
		}

		if locked, err := collection1.PropertyGetLocked(); err == nil {
			if !locked {
				t.Error("Expected collection1 'Lock' property to be true, got false")
			}
		} else {
			t.Error(err)
		}

		if locked, err := collection2.PropertyGetLocked(); err == nil {
			if locked {
				t.Errorf("Expected collection1 'Lock' property to be false, got true")
			}
		} else {
			t.Error(err)
		}

		if locked, err := collection3.PropertyGetLocked(); err == nil {
			if !locked {
				t.Errorf("Expected collection3 'Lock' property to be true, got false")
			}
		} else {
			t.Error(err)
		}

		expectedlocked := 0

		for _, lockedCollection := range locked {
			if lockedCollection == collection1.ObjectPath {
				expectedlocked++
			}
			if lockedCollection == collection3.ObjectPath {
				expectedlocked++
			}
		}

		if expectedlocked != 2 {
			t.Errorf("Expected 2 locked collections, got: %d", expectedlocked)
		}

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Locked <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Created >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READ UInt64 Created ;
	*/
	t.Run("Collection Property - Created", func(t *testing.T) {

		ssClient, _ := client.New()

		collection1, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")

		time.Sleep(time.Second * 2)

		collection2, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")

		collection1epoch, err := collection1.PropertyCreated()

		if err != nil {
			t.Error(err)
		}

		collection2epoch, err := collection2.PropertyCreated()

		if err != nil {
			t.Error(err)
		}

		if (collection2epoch - collection1epoch) < 2 {
			t.Errorf("collection2 'Created' property is not greater than collection1")
		}

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Created <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Modified >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READ UInt64 Modified ;
	*/
	t.Run("Collection Property - Modified", func(t *testing.T) {

		ssClient, _ := client.New()

		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "")

		collectionEpochBefore, err := collection.PropertyModified()

		if err != nil {
			t.Error(err)
		}

		time.Sleep(time.Second * 2)

		collection.PropertySetLabel("Modified-Test")

		collectionEpochAfter, err := collection.PropertyModified()

		if err != nil {
			t.Error(err)
		}

		if (collectionEpochAfter - collectionEpochBefore) < 2 {
			t.Errorf("collection 'Modified' property has not changes after set 'Label'")
		}

		// TODO: Other scenarios i.e. changing other properties

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Modified <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

}
