package service_test

import (
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

////////////////////////////// OpenSession //////////////////////////////

func Test_SecretServiceOpenSession(t *testing.T) {

	/*
		OpenSession ( IN String algorithm,
		IN Variant input,
		OUT Variant output,
		OUT String cookie);
	*/

	t.Run("dh-ietf1024-sha256-aes128-cbc-pkcs7 algorithm", func(t *testing.T) {

		ssClient, _ := client.New()
		err := ssClient.SecretServiceOpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("OpenSession failed. Error: %v", err)
		}

		if len(ssClient.CliSession.Cookie) != 32 {
			t.Errorf("Unexpected CLI cookie length. Expected 32, got '%d'",
				len(ssClient.CliSession.Cookie))
		}

	})

}
