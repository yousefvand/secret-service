package client_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
)

func Test_Item_Properties(t *testing.T) {

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Locked >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READ Boolean Locked ;
	*/
	t.Run("Collection Property - Locked", func(t *testing.T) {

		ssClient, _ := client.New()

		// open session
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session. Error: %v", err)
		}

		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "item property")

		serviceCollection := Service.GetCollectionByAlias("item property")

		////////////////////////////// item1 //////////////////////////////

		properties1 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("some item"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
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
			t.Errorf("CreateItem failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", itemPrompt)
		}

		if collection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at client side: %s", item1.ObjectPath)
		}

		if serviceCollection.GetItemByPath(item1.ObjectPath) == nil {
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

		iv2, cipherData2, err2 := client.AesCBCEncrypt([]byte("Victoria2"), session.SymmetricKey)

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

		if collection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at client side: %s", item2.ObjectPath)
		}

		if serviceCollection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at service side: %s", item2.ObjectPath)
		}

		////////////////////////////// Lock item1 //////////////////////////////

		locked, prompt, err := ssClient.Lock([]dbus.ObjectPath{item1.ObjectPath})

		if err != nil {
			t.Errorf("Service lock failed. Error: %v", err)
		}

		if prompt != "/" {
			t.Errorf("Service lock returned unknown prompt: %s", prompt)
		}

		////////////////////////////// Checks //////////////////////////////

		if len(locked) != 1 || locked[0] != item1.ObjectPath {
			t.Error("Wrong 'Lock' result")
		}

		if item1.Locked == false {
			t.Errorf("item1 is not locked at client side: %s", item1.ObjectPath)
		}

		if !Service.GetItemByPath(item1.ObjectPath).Locked {
			t.Errorf("item1 is not locked at service side: %s", item1.ObjectPath)
		}

		if item2.Locked == true {
			t.Errorf("item2 is locked at client side: %s", item2.ObjectPath)
		}

		if Service.GetItemByPath(item2.ObjectPath).Locked {
			t.Errorf("item2 is locked at service side: %s", item2.ObjectPath)
		}

		if locked, err := item1.PropertyGetLocked(); err == nil {
			if !locked {
				t.Error("Expected item1 'Lock' property to be true, got false")
			}
		} else {
			t.Error(err)
		}

		if locked, err := item2.PropertyGetLocked(); err == nil {
			if locked {
				t.Error("Expected item2 'Lock' property to be false, got true")
			}
		} else {
			t.Error(err)
		}

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Locked <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Attributes >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READWRITE Dict<String,String> Attributes ;
	*/
	t.Run("Collection Property - Attributes", func(t *testing.T) {

		ssClient, _ := client.New()

		// open session
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session. Error: %v", err)
		}

		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "item property")

		serviceCollection := Service.GetCollectionByAlias("item property")

		// Add item

		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("some item"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
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

		// Add first item
		item, prompt, err := collection.CreateItem(properties, secretApi, true)

		if err != nil {
			t.Errorf("CreateItem failed. Error: %v", err)
		}

		if prompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", prompt)
		}

		if collection.GetItemByPath(item.ObjectPath) == nil {
			t.Errorf("No such item1 at client side: %s", item.ObjectPath)
		}

		if serviceCollection.GetItemByPath(item.ObjectPath) == nil {
			t.Errorf("No such item1 at service side: %s", item.ObjectPath)
		}

		serviceItem := Service.GetItemByPath(item.ObjectPath)

		if !reflect.DeepEqual(serviceItem.LookupAttributes, item.LookupAttributes) {
			t.Errorf("Service side attributes: '%v' are different from client side: '%v'",
				serviceItem.LookupAttributes, item.LookupAttributes)
		}

		////////////////////////////// Set Attributes //////////////////////////////

		attributes := map[string]string{
			"w": "x",
			"y": "z",
		}
		err = item.PropertySetAttributes(attributes)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(item.LookupAttributes, attributes) {
			t.Errorf("Expected client side item attributes to be: '%v', got: '%v'",
				attributes, item.LookupAttributes)
		}

		if !reflect.DeepEqual(serviceItem.LookupAttributes, attributes) {
			t.Errorf("Expected service side item attributes to be: '%v', got: '%v'",
				attributes, serviceItem.LookupAttributes)
		}

		propertyAttributes, err := item.PropertyGetAttributes()

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(propertyAttributes, attributes) {
			t.Errorf("Expected item property attributes to be: '%v', got: '%v'",
				attributes, propertyAttributes)
		}

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Attributes <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Label >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READWRITE String Label ;
	*/
	t.Run("Collection Property - Label", func(t *testing.T) {

		ssClient, _ := client.New()

		// open session
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session. Error: %v", err)
		}

		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "item property")

		serviceCollection := Service.GetCollectionByAlias("item property")

		// Add item

		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("item-before"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
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

		// Add first item
		item, prompt, err := collection.CreateItem(properties, secretApi, true)

		if err != nil {
			t.Errorf("CreateItem failed. Error: %v", err)
		}

		if prompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", prompt)
		}

		if collection.GetItemByPath(item.ObjectPath) == nil {
			t.Errorf("No such item1 at client side: %s", item.ObjectPath)
		}

		if serviceCollection.GetItemByPath(item.ObjectPath) == nil {
			t.Errorf("No such item1 at service side: %s", item.ObjectPath)
		}

		serviceItem := Service.GetItemByPath(item.ObjectPath)

		if item.Label != "item-before" {
			t.Errorf("Expected item's client side label to be: 'item-before', got: '%v'", item.Label)
		}

		if serviceItem.Label != "item-before" {
			t.Errorf("Expected item's service side label to be: 'item-before', got: '%v'", serviceItem.Label)
		}

		itemLabel, err := item.PropertyGetLabel()

		if err != nil {
			t.Error(err)
		}

		if itemLabel != "item-before" {
			t.Errorf("Expected item's label property to be: 'item-before', got: '%v'", itemLabel)
		}

		////////////////////////////// Set item Label //////////////////////////////

		err = item.PropertySetLabel("item-after")

		if err != nil {
			t.Error(err)
		}

		if item.Label != "item-after" {
			t.Errorf("Expected item's client side label to be: 'item-after', got: '%v'", item.Label)
		}

		if serviceItem.Label != "item-after" {
			t.Errorf("Expected item's service side label to be: 'item-after', got: '%v'", serviceItem.Label)
		}

		itemLabel, err = item.PropertyGetLabel()

		if err != nil {
			t.Error(err)
		}

		if itemLabel != "item-after" {
			t.Errorf("Expected item's label property to be: 'item-after', got: '%v'", itemLabel)
		}

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Label <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Created >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READ UInt64 Created ;
	*/
	t.Run("Collection Property - Created", func(t *testing.T) {

		ssClient, _ := client.New()

		// open session
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session. Error: %v", err)
		}

		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "item property")

		serviceCollection := Service.GetCollectionByAlias("item property")

		////////////////////////////// item1 //////////////////////////////

		properties1 := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("some item"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
				"e": "f",
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
			t.Errorf("CreateItem failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", itemPrompt)
		}

		if collection.GetItemByPath(item1.ObjectPath) == nil {
			t.Errorf("No such item1 at client side: %s", item1.ObjectPath)
		}

		if serviceCollection.GetItemByPath(item1.ObjectPath) == nil {
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

		iv2, cipherData2, err2 := client.AesCBCEncrypt([]byte("Victoria2"), session.SymmetricKey)

		if err != nil {
			t.Errorf("encryption2 error: %v", err2)
		}

		secretApi2 := client.NewSecretApi()
		secretApi2.ContentType = "text/plain"
		secretApi2.Session = session.ObjectPath
		secretApi2.Parameters = iv2
		secretApi2.Value = cipherData2

		time.Sleep(time.Second * 2)

		// Add second item
		item2, itemPrompt, itemErr := collection.CreateItem(properties2, secretApi2, true)

		if itemErr != nil {
			t.Errorf("CreateItem failed. Error: %v", itemErr)
		}

		if itemPrompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", itemPrompt)
		}

		if collection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at client side: %s", item2.ObjectPath)
		}

		if serviceCollection.GetItemByPath(item2.ObjectPath) == nil {
			t.Errorf("No such item2 at service side: %s", item2.ObjectPath)
		}

		////////////////////////////// Compare Creation time //////////////////////////////

		item1epoch, err := item1.PropertyCreated()

		if err != nil {
			t.Error(err)
		}

		item2epoch, err := item2.PropertyCreated()

		if err != nil {
			t.Error(err)
		}

		if (item2epoch - item1epoch) < 2 {
			t.Errorf("item2 'Created' property is not greater than item1")
		}

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Created <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

	/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Modified >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

	/*
		READ UInt64 Modified ;
	*/
	t.Run("Collection Property - Modified", func(t *testing.T) {

		ssClient, _ := client.New()

		// open session
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("failed to open session. Error: %v", err)
		}

		collection, _, _ := ssClient.CreateCollection(map[string]dbus.Variant{}, "item property")

		serviceCollection := Service.GetCollectionByAlias("item property")

		// Add item

		properties := map[string]dbus.Variant{
			"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("item-before"),
			"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
				"a": "b",
				"c": "d",
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

		item, prompt, err := collection.CreateItem(properties, secretApi, true)

		serviceItem := serviceCollection.GetItemByPath(item.ObjectPath)

		if err != nil {
			t.Errorf("CreateItem failed. Error: %v", err)
		}

		if prompt != "/" {
			t.Errorf("wrong prompt for CreateItem: %v", prompt)
		}

		if collection.GetItemByPath(item.ObjectPath) == nil {
			t.Errorf("No such item1 at client side: %s", item.ObjectPath)
		}

		if serviceItem == nil {
			t.Errorf("No such item1 at service side: %s", item.ObjectPath)
		}

		////////////////////////////// Read Property Created //////////////////////////////

		itemEpochBefore, err := item.PropertyModified()

		if err != nil {
			t.Error(err)
		}

		time.Sleep(time.Second * 2)

		// Change item
		err = item.PropertySetLabel("new-label")

		if err != nil {
			t.Error(err)
		}

		itemEpochAfter, err := item.PropertyModified()

		if err != nil {
			t.Error(err)
		}

		if (itemEpochAfter - itemEpochBefore) < 2 {
			t.Errorf("item 'Modified' property has not changes after set 'Label'")
		}

		if serviceItem != nil && serviceItem.Modified != item.Modified {
			t.Errorf("Service Modified time: %d, Client Modified time: %d. Out of sync!",
				serviceItem.Modified, item.Modified)
		}

		// TODO: Other scenarios i.e. changing other properties

	})

	/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Modified <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

}
