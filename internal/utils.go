package internal

import "os"

// fileOrFolderExists returns whether the given file or directory exists
// Credit: https://stackoverflow.com/a/10510783
func fileOrFolderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
