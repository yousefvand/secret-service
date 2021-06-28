package internal

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type LogLevel uint8

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel LogLevel = iota // 0
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel // 1
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel // 2
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel // 3
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel // 4
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel // 5
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel // 6
)

type Config struct {
	// Config file version
	Version string `yaml:"version"`
	// Encrypt database using AES-CBC-256
	Encryption bool `yaml:"encryption"`
	// Desktop notification icon
	Icon string `yaml:"icon"`
	// Absolute path to log file
	LogFile string `yaml:"logFile"`
	// Logger is enabled or not
	Logging bool `yaml:"logging"`
	// Log file format: 'text' or 'json'
	LogFormat string `yaml:"logFormat"`
	// Level of logging verbosity
	// 0-6  where 6 is the most verbose and 0 logs very severe events
	LogLevel LogLevel `yaml:"logLevel"`
	// Maximum log size in MB before rotation
	LogMaxSize int `yaml:"logMaxSize"`
	// Maximum number of rotated log files (backups)
	LogMaxBackups int `yaml:"logMaxBackups"`
	// Maximum age (in days) of a log file before rotation
	// Between size and age which comes first makes log rotate
	LogMaxAge int `yaml:"logMaxAge"`
	// Compress backup log files (true) or not (false)
	LogCompress bool `yaml:"logCompress"`
	// Log report caller function (makes logs larger)
	LogReportCaller bool `yaml:"logReportCaller"`
}

// NewConfig returns a new instance of Config
func NewConfig() *Config {
	return &Config{}
}

// load configurations from file
func (config *Config) Load(app *AppData) {

	serviceHome := app.Service.Home
	filePath := filepath.Join(serviceHome, "config.yaml")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Warnf("Cannot open config file at: %s. Using default config.", filePath)
		// app.Notify("No config file", "Cannot find config file. Using default configurations.", time.Second*5)
		createDefaultConfig(serviceHome)
		data, _ = ioutil.ReadFile(filePath)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Warnf("found malformed config file: '%s'. Using default config.", filePath)
		app.Notify("Malformed config file", "Config file is malformed. Using default configurations.", time.Second*5)
		createDefaultConfig(serviceHome)
		data, _ = ioutil.ReadFile(filePath)
		_ = yaml.Unmarshal(data, config)
	}

	fillMissingConfigurations(config)
}

/* Just in case needed

func (config *Config) Save() error {

	filePath := GetConfigFilePath()
	data, errParsing := yaml.Marshal(config)
	if errParsing != nil {
		log.Errorf("Cannot marshal config data. Error: %v", errParsing)
		return errParsing
	}

	errWriting := ioutil.WriteFile(filePath, data, 0600)
	if errWriting != nil {
		log.Errorf("Cannot save config file at: '%s'. Error: %v", filePath, errWriting)
		return errWriting
	}

	return nil
}

*/

// fillMissingConfigurations fills missing important configurations with default values
func fillMissingConfigurations(config *Config) *Config {

	if config.Version == "" {
		config.Icon = "0.1.0"
	}

	if config.Icon == "" {
		config.Icon = "view-private"
	}

	if config.LogLevel > 6 {
		config.LogLevel = LogLevel(3) // Default: Warn
	}
	// TODO: Complete
	if exist, err := fileOrFolderExists(filepath.Dir(config.LogFile)); err != nil || !exist {
		config.LogFile = ""
	}

	if logFormat := strings.ToLower(config.LogFormat); logFormat != "text" && logFormat != "json" {
		config.LogFormat = "text"
	}

	if config.LogMaxSize < 1 {
		config.LogMaxSize = 1
	}

	if config.LogMaxBackups < 1 {
		config.LogMaxBackups = 1
	}

	if config.LogMaxAge < 1 {
		config.LogMaxAge = 1
	}

	return config
}

// createDefaultConfig creates default config file in case it is missing
func createDefaultConfig(appHome string) {

	configFilePath := filepath.Join(appHome, "config.yaml")

	errWriteConfig := ioutil.WriteFile(configFilePath, defaultConfig, 0600)

	if errWriteConfig != nil {
		log.Panicf("Writing default configurations to '%s' failed. Error: %v",
			configFilePath, errWriteConfig)
	}

}

// template for configuration file with default values
var defaultConfig []byte = []byte(`# Config file version
version: '0.1.0'

# Encrypt database using AES-CBC-256
# You need to set a MASTERPASSWORD in '/etc/systemd/user/secretserviced.service'
# File with EXACTLY 32 characters length or this configuration is ignored
encryption: true

# A system icon as string i.e. "flag" used in notifications
icon: 'view-private'

# 0-6  where 6 is the most verbose and 0 logs very severe events
# 0: Panic, 1: Fatal, 2: Error, 3: Warning, 4: Info, 5: Debug, 6: Trace
logLevel: 6

# Absolute path to log file
logFile: ''

# Enable logger
logging: true

# Log file format: 'text' or 'json'
logFormat: 'text'

# Maximum log size in MB
logMaxSize: 5

# Maximum number of log file backups
logMaxBackups: 10

# Maximum age (in days) of a log file
logMaxAge: 30

# Compress backup log files (true) or not (false)
logCompress: true

# Log report caller function
logReportCaller: false
`)
