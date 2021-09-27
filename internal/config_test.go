package internal

import (
	"testing"
)

// const configFileName = "config.yaml"

func Test_Config(t *testing.T) {

	config := NewConfig()

	t.Run("NewConfig", func(t *testing.T) {

		if config == nil {
			t.Error("new config is null")
		}
	})

	t.Run("LoadConfig", func(t *testing.T) {

		app := NewApp()
		app.Service.Config.Home = tempDir(t)
		config.Load(app)

		if config.Icon == "" || config.Version == "" {
			t.Error("config has no icon")
		}
	})
}
