package client_test

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/yousefvand/secret-service/pkg/client"
	"github.com/yousefvand/secret-service/pkg/crypto"
)

/*
GetSecrets ( IN Array<ObjectPath> items,
             IN ObjectPath session,
             OUT Dict<ObjectPath,Secret> secrets);
*/

func TestClient_GetSecrets(t *testing.T) {

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
