package client_test

import (
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

/*
	OpenSession ( IN String algorithm,
	              IN Variant input,
	              OUT Variant output,
	              OUT ObjectPath result);
*/

func TestClient_OpenSession(t *testing.T) {

	t.Run("Service OpenSession - plain algorithm", func(t *testing.T) {

		ssClient, _ := client.New()

		session, err := ssClient.OpenSession(client.Plain)

		if err != nil {
			t.Errorf("OpenSession failed. Error: %v", err)
		}

		if session != nil && session.PublicKey != nil {
			t.Errorf("Unexpected public key. Expected null public key for plain algorithm")
		}

		// i.e. "/org/freedesktop/secrets/session/uuid..." uuid is 32 character
		if session != nil && session.ObjectPath[:33] != "/org/freedesktop/secrets/session/" {
			t.Errorf("invalid session path (should start with '/org/freedesktop/secrets/session/'): %s", session.ObjectPath)
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

	t.Run("Service OpenSession - dh-ietf1024-sha256-aes128-cbc-pkcs7 algorithm", func(t *testing.T) {

		ssClient, _ := client.New()
		session, err := ssClient.OpenSession(client.Dh_ietf1024_sha256_aes128_cbc_pkcs7)

		if err != nil {
			t.Errorf("OpenSession failed. Error: %v", err)
		}

		if len(session.ServicePublicKey) != 128 {
			t.Errorf("Unexpected service public key length. Expected 128, got '%d'",
				len(session.ServicePublicKey))
		}

		// i.e. "/org/freedesktop/secrets/session/uuid..." uuid is 32 character
		if session.ObjectPath[:33] != "/org/freedesktop/secrets/session/" {
			t.Errorf("invalid session path (should start with '/org/freedesktop/secrets/session/'): %s", session.ObjectPath)
		}

		if len(session.ObjectPath) != 65 {
			t.Errorf("invalid objectPath (length is not 65): %s", session.ObjectPath)
		}

		if !ssClient.HasSession(session.ObjectPath) {
			t.Errorf("session doesn't exist at client side: %s", session.ObjectPath)
		}

		if !Service.HasSession(session.ObjectPath) {
			t.Errorf("session doesn't exist at service side: %s", session.ObjectPath)
		}

	})

	t.Run("Service OpenSession - unsupported algorithm", func(t *testing.T) {

		ssClient, _ := client.New()
		session, err := ssClient.OpenSession(client.Unsupported)

		if err == nil {
			t.Errorf("OpenSession unsupported algorithm raised no error. Error: %v", err)
		}

		if session != nil && session.PublicKey != nil {
			t.Errorf("Unexpected public key. Expected null for unsupported algorithm")
		}

	})

}
