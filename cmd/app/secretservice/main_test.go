package main

import (
	"testing"
)

func Test_main(t *testing.T) {

	t.Run("Cli", func(t *testing.T) {
		main()
		t.Log("Cli app implemented")
	})
}
