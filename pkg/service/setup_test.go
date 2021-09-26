package service_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yousefvand/secret-service/pkg/service"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Service *service.Service

// TestMain runs before all other tests. It starts service, waits for
// it to be up then runs other tests and finally shutdowns the service.
func TestMain(m *testing.M) {

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "../../logs/service-test.log", // Log file relative path
		MaxSize:    5,                             // MB
		MaxBackups: 1,
		MaxAge:     30,    // days
		Compress:   false, // disabled by default
	}
	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = time.RFC1123Z // RFC3339 or "02-01-2006 15:04:05"
	logFormatter.FullTimestamp = true
	log.SetFormatter(logFormatter)
	log.SetLevel(log.Level(log.TraceLevel))
	log.SetOutput(lumberjackLogger)

	Service = service.New()
	Service.Config.Home, _ = ioutil.TempDir("", "secret-service")
	ctx, cancel := context.WithCancel(context.Background())
	go Service.Start(ctx) // start secret service

	<-Service.ServiceReadyChan    // wait for service to be up and ready
	errCode := m.Run()            // run other tests and get the error code if any
	cancel()                      // shutdown secret service
	<-Service.ServiceShutdownChan // wait for service to signal shutting down
	os.Exit(errCode)              // errCode = 0 means no error in tests
}
