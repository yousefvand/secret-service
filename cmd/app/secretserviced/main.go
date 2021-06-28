// secretserviced man package (entry point)
// the actual work starts at app package
package main

import (
	"context"

	"github.com/yousefvand/secret-service/cmd/app"
)

// Entry point for secretserviced daemon
func main() {
	// TODO: Argument parsing

	// Run the service with empty context
	// context is also used in unit tests
	app.Run(context.Background())
}
