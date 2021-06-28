// secretserviced man package
// the actual work starts at app package

package main

import (
	"testing"

	"github.com/yousefvand/secret-service/cmd/app"
)

func Test_main(t *testing.T) {

	t.Run("secretserviced - entry point", func(t *testing.T) {
		go main()
		<-app.App.Service.ServiceReadyChan
		t.Log("service is up and ready")
	})
}
