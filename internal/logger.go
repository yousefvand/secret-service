package internal

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// SetupLogger sets up logger according to config file
func SetupLogger(app *AppData) {

	if !app.Config.Logging { // logging is disabled
		log.SetLevel(log.PanicLevel) // set log level to least messages
		log.SetOutput(io.Discard)    // discard output
		return
	}

	// No log file specified, use default:
	// ~/.secret-service/secretserviced/logs/secretserviced.log
	if app.Config.LogFile == "" {
		app.Config.LogFile = filepath.Join(app.Service.Home, "logs", "secretserviced.log")
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   app.Config.LogFile,       // Log file abbsolute path
		MaxSize:    app.Config.LogMaxSize,    // Log file size in MB
		MaxBackups: app.Config.LogMaxBackups, // Maximum of log file backups
		MaxAge:     app.Config.LogMaxAge,     // Log file age in days
		Compress:   app.Config.LogCompress,   // disabled by default
	}

	// Fork writing into two outputs (log file and stderr)
	// So writing to multiWriter results in writing to both
	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	// Log file in TEXT format
	logTextFormatter := new(log.TextFormatter)
	logTextFormatter.TimestampFormat = time.RFC1123Z // RFC3339 or "02-01-2006 15:04:05"
	logTextFormatter.FullTimestamp = true

	// Log file in JSON format
	logJSONFormatter := new(log.JSONFormatter)
	logJSONFormatter.TimestampFormat = time.RFC1123Z // RFC3339 or "02-01-2006 15:04:05"
	logJSONFormatter.DisableHTMLEscape = true

	if strings.ToLower(app.Config.LogFormat) == "json" {
		log.SetFormatter(logJSONFormatter)
	} else { // if not JSON use TEXT
		log.SetFormatter(logTextFormatter)
	}

	log.SetLevel(log.Level(app.Config.LogLevel))

	log.SetReportCaller(app.Config.LogReportCaller)

	if os.Getenv("ENV") == "TEST" { // Don't mess up test output
		log.SetOutput(lumberjackLogger)
	} else {
		log.SetOutput(multiWriter)
	}

}
