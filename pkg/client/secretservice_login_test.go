package client_test

import (
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

/*
	Login ( IN  String serialnumber,
	        IN  Array<Byte> passwordhash,
					IN  Array<Byte> passwordhash_iv,
					OUT Array<Byte> cookie,
					OUT Array<Byte> cookie_iv
					OUT String result);
*/

func TestClient_SecretServiceLogin(t *testing.T) {

	t.Run("Login - failure", func(t *testing.T) {

		ssClient, _ := client.New()
		err := ssClient.SecretServiceCreateSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("CreateSession failed. Error: %v", err)
		}

		// write a password file
		secret := "Alpha" + "Bravo"
		hasher := sha512.New()
		hasher.Write([]byte(secret))
		hash := hex.EncodeToString(hasher.Sum(nil))

		err = Service.WritePasswordFile(hash)

		if err != nil {
			t.Errorf("Cannot write password file. Error: %v", err)
		}

		// try login
		secret = "Charlie" + "Tango"
		hasher = sha512.New()
		hasher.Write([]byte(secret))
		hash = hex.EncodeToString(hasher.Sum(nil))

		iv, cipher, err := client.AesCBCEncrypt([]byte(hash),
			ssClient.SecretService.Session.SymmetricKey)

		if err != nil {
			t.Errorf("Encryption failed. Error: %v", err)
		}

		encryptedCookie, cookie_iv, result, err := ssClient.SecretServiceLogin(
			ssClient.SecretService.Session.SerialNumber, cipher, iv)

		if err != nil {
			t.Errorf("Login returned error: %v", err)
		}

		if result != "wrong password" {
			t.Errorf("Expected 'wrong password' got: %v", result)
		}

		_, err = client.AesCBCDecrypt(cookie_iv, encryptedCookie, ssClient.SecretService.Session.SymmetricKey)

		if err == nil {
			t.Error("Decryption didn't fail!")
		}
	})
}
