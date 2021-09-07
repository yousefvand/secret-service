package client_test

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

/*
	SetPassword ( IN  String      serialnumber
		            IN  Array<Byte> oldpassword,
								IN  Array<Byte> oldpassword_iv,
								IN  Array<Byte> newpassword,
								IN  Array<Byte> newpassword_iv,
								IN  Array<Byte> oldSalt,
								IN  Array<Byte> oldSalt_iv,
								IN  Array<Byte> newSalt,
								IN  Array<Byte> newSalt_iv
								OUT String result);
*/

func TestClient_SecretServiceSetPassword(t *testing.T) {

	t.Run("SetPassword - empty", func(t *testing.T) {

		ssClient, _ := client.New()
		err := ssClient.SecretServiceCreateSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("CreateSession failed. Error: %v", err)
		}

		if result := bytes.Compare(Service.SecretService.Session.SymmetricKey, ssClient.SecretService.Session.SymmetricKey); result != 0 {
			t.Errorf("Symmetric keys are not equal!")
		}

		////////////////////////////// SetPassword (empty) //////////////////////////////

		oldPassword_iv, oldPassword_cipher, _ := client.AesCBCEncrypt([]byte(""), ssClient.SecretService.Session.SymmetricKey)
		newPassword_iv, newPassword_cipher, _ := client.AesCBCEncrypt([]byte("Victoria"), ssClient.SecretService.Session.SymmetricKey)
		oldSalt_iv, oldSaltCipher, _ := client.AesCBCEncrypt([]byte("Salt"), ssClient.SecretService.Session.SymmetricKey)
		newSalt_iv, newSaltCipher, _ := client.AesCBCEncrypt([]byte("Salt"), ssClient.SecretService.Session.SymmetricKey)

		result, err := ssClient.SecretServiceSetPassword(ssClient.SecretService.Session.SerialNumber,
			oldPassword_cipher, oldPassword_iv, newPassword_cipher, newPassword_iv,
			oldSaltCipher, oldSalt_iv, newSaltCipher, newSalt_iv)

		if err != nil {
			t.Errorf("SetPassword Failed.Error: %v", err)
		}

		if result != "ok" {
			t.Errorf("Expected 'ok' got: %s", result)
		}

		t.Run("SetPassword - change", func(t *testing.T) {

			ssClient, _ := client.New()
			err := ssClient.SecretServiceCreateSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

			if err != nil {
				t.Errorf("CreateSession failed. Error: %v", err)
			}

			secret := "OldSalt" + "OldVictoria"
			hasher := sha512.New()
			hasher.Write([]byte(secret))
			hash := hex.EncodeToString(hasher.Sum(nil))

			err = Service.WritePasswordFile(hash)

			if err != nil {
				t.Errorf("Cannot write password file. Error: %v", err)
			}

			oldPassword_iv, oldPassword_cipher, _ := client.AesCBCEncrypt([]byte("OldVictoria"), ssClient.SecretService.Session.SymmetricKey)
			newPassword_iv, newPassword_cipher, _ := client.AesCBCEncrypt([]byte("NewVictoria"), ssClient.SecretService.Session.SymmetricKey)
			oldSalt_iv, oldSaltCipher, _ := client.AesCBCEncrypt([]byte("OldSalt"), ssClient.SecretService.Session.SymmetricKey)
			newSalt_iv, newSaltCipher, _ := client.AesCBCEncrypt([]byte("NewSalt"), ssClient.SecretService.Session.SymmetricKey)

			result, err := ssClient.SecretServiceSetPassword(ssClient.SecretService.Session.SerialNumber,
				oldPassword_cipher, oldPassword_iv, newPassword_cipher, newPassword_iv,
				oldSaltCipher, oldSalt_iv, newSaltCipher, newSalt_iv)

			if err != nil {
				t.Errorf("SetPassword Failed.Error: %v", err)
			}

			if result != "ok" {
				t.Errorf("Expected 'ok' got: %s", result)
			}

		})
	})
}
