package service_test

import (
	"testing"

	"github.com/yousefvand/secret-service/pkg/client"
)

////////////////////////////// CreateSession //////////////////////////////

func Test_SecretServiceCommand(t *testing.T) {

	t.Run("ping", func(t *testing.T) {

		ssClient, _ := client.New()
		response, _ := ssClient.SecretServiceCommand("ping", "")

		if response != "pong" {
			t.Errorf("Expected 'pong' got: %s", response)
		}

		response, _ = ssClient.SecretServiceCommand("foo", "")

		if response != "bar" {
			t.Errorf("Expected 'bar' got: %s", response)
		}

		response, _ = ssClient.SecretServiceCommand("baz", "")

		if response != "unknown" {
			t.Errorf("Expected 'unknown' got: %s", response)
		}

	})

}
