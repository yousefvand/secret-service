package service_test

import (
	"bytes"
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

////////////////////////////// CreateSession //////////////////////////////

func Test_SecretServiceCreateSession(t *testing.T) {

	/*
		CreateSession ( IN String algorithm,
		                IN Variant input,
		                OUT Variant output,
		                OUT String serialnumber);
	*/

	t.Run("dh-ietf1024-sha256-aes128-cbc-pkcs7 algorithm", func(t *testing.T) {

		ssClient, _ := client.New()
		err := ssClient.SecretServiceCreateSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("CreateSession failed. Error: %v", err)
		}

		if len(ssClient.SecretService.Session.SerialNumber) != 32 {
			t.Errorf("Unexpected CLI serialnumber length. Expected 32, got '%d'",
				len(ssClient.SecretService.Session.SerialNumber))
		}

		if result := bytes.Compare(Service.SecretService.Session.SymmetricKey, ssClient.SecretService.Session.SymmetricKey); result != 0 {
			t.Errorf("Symmetric keys are not equal!")
		}

	})

}

////////////////////////////// SetPassword //////////////////////////////

func Test_SecretServiceSetPassword(t *testing.T) {

	/*
		SetPassword ( IN  String      serialnumber
			            IN  Array<Byte> oldpassword,
									IN  Array<Byte> oldpassword_iv,
									IN  Array<Byte> newpassword,
									IN  Array<Byte> newpassword_iv,
									IN  Array<Byte> salt,
									IN  Array<Byte> salt_iv
									OUT String result);
	*/

	t.Run("SetPassword", func(t *testing.T) {

	})
}
