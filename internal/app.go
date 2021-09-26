// Internal operations mostly related
// to App, configurations, and logger
package internal

import (
	"time"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"github.com/yousefvand/secret-service/pkg/service"
)

// App data structure, all other data are
// children of App in a hierarchial model
type AppData struct {
	// Service (secretserviced) home directory
	ServiceHome string
	// dbus session connection
	Connection *dbus.Conn
	// configurations data structure
	Config *Config
	// Service data structure
	Service *service.Service
}

// NewApp creates and initialize a new instance of App
func NewApp() *AppData {

	// create a new App instance
	app := &AppData{}

	// create a new Config instance
	app.Config = NewConfig()

	// create a new Service instance
	app.Service = service.New()

	// Create secret-service related directories and return secretserviced
	//  home directory (~/.secret-service/secretserviced/)
	// .secret-service
	// ├── secretservice-cli
	// │   └── logs
	// └── secretserviced
	//     ├── config.yaml
	//     ├── db.json
	//     └── logs
	//         └── secretserviced.log
	app.Service.Config.Home = SetupHomeDirectories()

	// Create a dbus session connection (singleton)
	connection, err := dbus.SessionBus()
	if err != nil {
		log.Panicf("Cannot connect to session dbus. Error: %v", err)
	}
	app.Connection = connection
	app.Service.Connection = connection

	return app
}

// Load configurations and setup logger
func (app *AppData) Load() {
	app.Config.Load(app)
	app.Service.Config.AllowDbExport = app.Config.AllowDbExport
	app.Service.Config.EncryptDatabase = app.Config.Encryption
	app.SetupLogger()
}

// SetupLogger sets up logger based on configurations
func (app *AppData) SetupLogger() {
	SetupLogger(app)
}

// Notify sends a desktop notification
func (app *AppData) Notify(title string, body string, duration time.Duration) {

	icon := app.Config.Icon

	if icon == "" {
		icon = "view-private" // or "flag"
	}

	Notify(app.Connection, "Secret Service", icon, title, body, duration)
}

// Notify sends a desktop notifications (standalone version)
func Notify(connection *dbus.Conn, appName string, icon string, title string, body string, duration time.Duration) error {

	obj := connection.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, appName, uint32(0),
		icon, title, body, []string{},
		map[string]dbus.Variant{}, int32(duration.Milliseconds()))
	if call.Err != nil {
		log.Errorf("Calling dbus notification failed. Error: %v", call.Err)
		return call.Err
	}
	return nil
}
