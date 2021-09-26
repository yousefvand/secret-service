package service_test

import (
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func Test_ReadPasswordFile(t *testing.T) {

	// t.Run("Read 'password.yaml' file - no file", func(t *testing.T) {

	// 	password := Service.ReadPasswordFile()
	// 	if password != "" {
	// 		t.Errorf("Expected no password, got: %s", password)
	// 	}

	// })

	t.Run("Read 'password.yaml' file - empty hash", func(t *testing.T) {

		// write empty hash
		passwordFile := filepath.Join(Service.Config.Home, "password.yaml")
		errWritePasswordFile := ioutil.WriteFile(passwordFile, emptyPasswordFile, 0600)

		if errWritePasswordFile != nil {
			t.Errorf("Cannot write password file. Error: %v", errWritePasswordFile)
		}

		password := Service.ReadPasswordFile()
		if password != "" {
			t.Errorf("Expected no password, got: %s", password)
		}

	})

	t.Run("Read 'password.yaml' file - some hash", func(t *testing.T) {

		// write some hash
		passwordFile := filepath.Join(Service.Config.Home, "password.yaml")
		errWritePasswordFile := ioutil.WriteFile(passwordFile, fullPasswordFile, 0600)

		if errWritePasswordFile != nil {
			t.Errorf("Cannot write password file. Error: %v", errWritePasswordFile)
		}

		password := Service.ReadPasswordFile()

		if len(password) < 1 {
			t.Error("Expected some password, got nothing!")
		}

	})

}

func Test_WritePasswordFile(t *testing.T) {

	t.Run("Write 'password.yaml' file - version & password hash", func(t *testing.T) {

		secret := "Salt" + "Victoria"
		hasher := sha512.New()
		hasher.Write([]byte(secret))
		hash := hex.EncodeToString(hasher.Sum(nil))

		err := Service.WritePasswordFile(hash)

		if err != nil {
			t.Errorf("Cannot write password file. Error: %v", err)
		}

		readHash := Service.ReadPasswordFile()

		if readHash != hash {
			t.Errorf("Expected %s password hash, got %s", hash, readHash)
		}
	})
}

var fullPasswordFile []byte = []byte(`# Password file version
version: 0.1.0

# Password hash: sha512(salt+password)
passwordHash: '3eb2a19ed384c94ba5a3f67e0bef39860dda717d0ca07d0fb4a49bb2e9db22f797cf8b68f47c6a9aa33fabaf561b9c4beb481718d06fd60558f8717c05db8022'
`)

var emptyPasswordFile []byte = []byte(`# Password file version
version: 0.1.0

# Password hash: sha512(salt+password)
passwordHash: ''
`)
