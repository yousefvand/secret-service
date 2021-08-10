package client_test

import (
	"bytes"
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

/*
	OpenSession ( IN String algorithm,
	              IN Variant input,
	              OUT Variant output,
	              OUT String serialnumber);
*/

func TestClient_SecretServiceOpenSession(t *testing.T) {

	t.Run("SecretService OpenSession - dh-ietf1024-sha256-aes128-cbc-pkcs7 algorithm", func(t *testing.T) {

		ssClient, _ := client.New()
		err := ssClient.SecretServiceOpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("OpenSession failed. Error: %v", err)
		}

		if len(ssClient.CliSession.SerialNumber) != 32 {
			t.Errorf("Unexpected CLI serialnumber length. Expected 32, got '%d'",
				len(ssClient.CliSession.SerialNumber))
		}

		if result := bytes.Compare(Service.CliSession.SymmetricKey, ssClient.CliSession.SymmetricKey); result != 0 {
			t.Errorf("Symmetric keys are not equal!")
		}

	})

}
