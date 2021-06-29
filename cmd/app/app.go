// managing logging and configuration loading
// and starting secretserviced daemon besides
// handling os signals.
package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/yousefvand/secret-service/internal"
)

/*

Exit codes:
1) SIGINT
2) SIGTERM
3) SIGQUIT
4) Non of above signal (except SIGHUP)
5) name "org/freedesktop/secrets" is already taken on dbus

*/

/* TODO: Remove or implement
type Commands uint8

const (
	Init Commands = iota
	Start
	Stop
	Reload
	Restart
	Clean
)
*/

// Main App struct and methods
// all other data are children
// of App in a hierarchy model
var App *internal.AppData

// init creates App instance and loads
// configurations and setup app logger
func init() {
	// A new instance of App with a "Service" as its child
	App = internal.NewApp()

	// Load configurations and setup logger
	App.Load()
}

// Run runs the App which means reading, (creating if not present) config
// file and creation of a dbus session connection and running the service
// (secretserviced) and storing data in a database, upon any change. Also
// handles OS signals (reloading/quitting). Context is used in unit tests.
func Run(ctx context.Context) {

	ctx, cancel := context.WithCancel(ctx)

	// Run secretserviced on a separate Goroutine
	go App.Service.Start(ctx)

	/* ========== Signal handling ========== */

	// Register for specific OS signals we are interested in
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// Channel for communicating exit signal
	exitChan := make(chan int)

	// OS signal handling
	go func(cancel func()) {
		for {
			signal := <-signalChan
			switch signal {
			case syscall.SIGHUP:
				log.Info("Received 'SIGHUP' signal. Reloading configurations...")
				App.Load()

			case syscall.SIGINT: // CTRL+C
				log.Info("***** Received 'SIGINT' signal. Exiting... *****")
				cancel()
				<-App.Service.ServiceShutdownChan
				exitChan <- 1

			case syscall.SIGTERM:
				log.Info("***** Received 'SIGTERM' signal. Exiting... *****")
				cancel()
				<-App.Service.ServiceShutdownChan
				exitChan <- 2

			case syscall.SIGQUIT:
				log.Info("***** Received 'SIGQUIT' signal. Exiting... *****")
				cancel()
				<-App.Service.ServiceShutdownChan
				exitChan <- 3

			default:
				log.Info("***** Received unknown signal. Exiting... *****")
				cancel()
				<-App.Service.ServiceShutdownChan
				exitChan <- 4
			}
		}
	}(cancel)

	os.Exit(<-exitChan) // exit secretserviced
}
