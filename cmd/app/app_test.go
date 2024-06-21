// managing logging and configuration loading
// and data saving under app name and starting
// secretserviced daemon also handling os signals.
package app

import (
	"context"
	"testing"
)

func TestRun(t *testing.T) {

	t.Run("App Run", func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())
		go Run(ctx)
		<-App.Service.ServiceReadyChan
		t.Log("service is up and ready")
		cancel()
		<-App.Service.ServiceShutdownChan
		t.Log("service is shutdown")
	})
	// TODO: Test OS signals
}
