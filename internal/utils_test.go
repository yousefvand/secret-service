package internal

import (
	"os"
	"path"
	"testing"
)

func Test_fileOrFolderExists(t *testing.T) {
	tmpFile := tempFile(t)
	fileExist, errTempFile := fileOrFolderExists(tmpFile)
	if errTempFile != nil {
		t.Errorf("File existence check failed. Error: %v", errTempFile)
	}
	if !fileExist {
		t.Errorf("File existence check failed.")
	}

	fakeFileExist, _ := fileOrFolderExists(path.Join(os.TempDir(), "secret-service.txt"))
	if fakeFileExist {
		t.Errorf("File existence check failed. No such a file!")
	}

	tmpDir := tempDir(t)
	dirExist, errTempDir := fileOrFolderExists(tmpDir)
	if errTempDir != nil {
		t.Errorf("Directory existence check failed. Error: %v", errTempDir)
	}
	if !dirExist {
		t.Errorf("Directory existence check failed.")
	}

	fakeDirExist, _ := fileOrFolderExists(path.Join(os.TempDir(), "secret-service"))
	if fakeDirExist {
		t.Errorf("Directory existence check failed. No such a directory!")
	}
}

func tempFile(t *testing.T) string {
	file, err := os.CreateTemp("", "secret-service.*.test")
	if err != nil {
		t.Errorf("Temporary file creation failed")
	}

	return file.Name()
}

func tempDir(t *testing.T) string {
	name, err := os.MkdirTemp("", "secret-service-*")
	if err != nil {
		t.Errorf("Temporary directory creation failed")
	}

	return name
}
