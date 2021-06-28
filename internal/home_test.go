package internal

import (
	"testing"
)

func Test_serviceHomeDirectory(t *testing.T) {

	// userHome, _ := os.UserHomeDir()
	// want := filepath.Join(userHome, ".config", "secretserviced")

	// t.Run("serviceHomeDirectory", func(t *testing.T) {
	// 	got, err := SetupHomeDirectories()
	// 	if err != nil {
	// 		t.Errorf("serviceHomeDirectory() error = %v", err)
	// 		return
	// 	}
	// 	if got != want {
	// 		t.Errorf("serviceHomeDirectory() = %v, want %v", got, want)
	// 	}
	// 	exists, errCreation := fileOrFolderExists(want)
	// 	if !exists || errCreation != nil {
	// 		t.Errorf("Service home directory not created. Error: %v", errCreation)
	// 	}
	// })
}
