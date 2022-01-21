// client package for 'secret service' as described at:
// http://standards.freedesktop.org/secret-service
package client

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
)

// New returns a new client connected to session dbus
func New() (*Client, error) {

	client := new(Client)
	client.SignalChan = make(chan *dbus.Signal)
	client.SessionsMutex = new(sync.RWMutex)
	client.CollectionsMutex = new(sync.RWMutex)
	client.Sessions = make(map[string]*Session)
	client.SecretService = &SecretService{}
	client.SecretService.Session = &SecretServiceCLiSession{}
	client.SecretService.Parent = client
	client.Collections = make(map[string]*Collection)

	connection, err := dbus.SessionBus()
	if err != nil {
		return nil, errors.New("cannot connect to session dbus. Error: " + err.Error())
	}
	client.Connection = connection

	// Watch for org.freedesktop.Secret.Service signals
	err = client.Connection.AddMatchSignal(
		dbus.WithMatchObjectPath("/org/freedesktop/secrets"),
		dbus.WithMatchInterface("org.freedesktop.Secret.Service"),
		dbus.WithMatchSender("org.freedesktop.secrets"),
	)

	if err != nil {
		return nil, errors.New("cannot watch for signals. Error: " + err.Error())
	}

	client.Connection.Signal(client.SignalChan)

	return client, nil
}

// Connected returns true if client is
// connected to session dbus otherwise false
func (client *Client) Connected() bool {
	return client.Connection.Connected()
}

// Disconnect from session dbus
// CAUTION: connection is shared by all clients
// by closing it all clients fail on subsequent operations
func (client *Client) Disconnect() error {
	err := client.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}

// Call performs low-level method Call on org.freedesktop.secrets objects
// Don't use this method directly unless you know what you are doing!
func (client *Client) Call(destination string, dbusPath dbus.ObjectPath,
	dbusInterface string, methodName string, args ...interface{}) (*dbus.Call, error) {

	dbusPath = dbus.ObjectPath(strings.TrimSpace(string(dbusPath)))
	dbusInterface = strings.TrimSpace(dbusInterface)

	if dbusPath[0] != '/' || dbusPath[len(dbusPath)-1:] == "/" {
		return nil, errors.New("Invalid dbusPath: " + string(dbusPath))
	}

	if dbusInterface[0] == '.' || dbusInterface[len(dbusInterface)-1:] == "." ||
		strings.Count(dbusInterface, ".") < 2 {
		return nil, errors.New("Invalid dbusInterface: " + dbusInterface)
	}

	client.DbusObject = client.Connection.Object(destination,
		dbus.ObjectPath(dbusPath))
	response := client.DbusObject.Call(dbusInterface+"."+methodName, 0, args...)
	return response, nil
}

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Session >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// AddSession adds a new session to client's session map
func (client *Client) AddSession(session *Session) {
	client.SessionsMutex.Lock()
	client.Sessions[string(session.ObjectPath)] = session
	client.SessionsMutex.Unlock()
}

// Remove remove a session from client's Sessions map
func (session *Session) Remove() error {
	client := session.Parent
	client.SessionsMutex.Lock()
	defer client.SessionsMutex.Unlock()
	_, ok := client.Sessions[string(session.ObjectPath)]
	if !ok {
		return errors.New("Session doesn't exist to be removed: " + string(session.ObjectPath))
	}
	delete(client.Sessions, string(session.ObjectPath))
	return nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Session <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Collection >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// AddCollection adds a new collection to client's collection map
func (client *Client) AddCollection(collection *Collection) {
	client.CollectionsMutex.Lock()
	defer client.CollectionsMutex.Unlock()
	client.Collections[string(collection.ObjectPath)] = collection
}

// RemoveCollection removes a collection from client's Collections map
func (client *Client) RemoveCollection(collection *Collection) error {
	client.CollectionsMutex.Lock()
	defer client.CollectionsMutex.Unlock()
	_, ok := client.Collections[string(collection.ObjectPath)]
	if !ok {
		return errors.New("Collection doesn't exist to be removed: " + string(collection.ObjectPath))
	}
	delete(client.Collections, string(collection.ObjectPath))
	return nil
}

// HasCollection returns true if collection exists otherwise false
func (client *Client) HasCollection(collectionPath dbus.ObjectPath) bool {
	client.CollectionsMutex.RLock()
	defer client.CollectionsMutex.RUnlock()
	_, ok := client.Collections[string(collectionPath)]
	return ok
}

// GetCollectionByPath returns a collection based on its path, otherwise null
func (client *Client) GetCollectionByPath(collectionPath dbus.ObjectPath) *Collection {
	client.CollectionsMutex.RLock()
	defer client.CollectionsMutex.RUnlock()

	for _, collection := range client.Collections {
		if collection.ObjectPath == collectionPath {
			return collection
		}
	}
	return nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Collection <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Signals >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// org.freedesktop.Secret.Service signals
type ServiceSignal uint8

const (
	CollectionCreated ServiceSignal = iota
	CollectionDeleted
	CollectionChanged
)

// WatchSignal watches for desired signal within a time period
// If signal is received it returns true, otherwise false
func (client *Client) WatchSignal(signal ServiceSignal, timeout ...time.Duration) (bool, error) {

	var signalName string

	signalTimeout := time.Second // default timeout
	if len(timeout) > 0 {
		signalTimeout = timeout[0]
	}

	switch signal {
	case CollectionCreated:
		signalName = "CollectionCreated"
	case CollectionDeleted:
		signalName = "CollectionDeleted"
	case CollectionChanged:
		signalName = "CollectionChanged"
	}

	select {
	case signal := <-client.SignalChan:
		if signal.Name == "org.freedesktop.Secret.Service."+signalName {
			return true, nil
		} else {
			return false, fmt.Errorf("expected 'org.freedesktop.Secret.Service.%s' signal got: %s", signalName, signal.Name)
		}
	case <-time.After(signalTimeout):
		return false, fmt.Errorf("receiving '%s' signal timed out", signalName)
	}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Signals <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Properties >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// PropertyGetCollections returns Collections property of Service
func (client *Client) PropertyGetCollections() ([]string, error) {

	busObject := client.Connection.Object("org.freedesktop.secrets",
		"/org/freedesktop/secrets")

	variant, err := busObject.GetProperty("org.freedesktop.Secret.Service.Collections")

	if err != nil {
		return nil, errors.New("Error getting property 'Collections': " + err.Error())
	}

	collectionspaths, ok := variant.Value().([]dbus.ObjectPath)

	if !ok {
		return nil, fmt.Errorf("invalid 'Collections' property type. Expected '[]dbus.ObjectPath', got '%T'", variant.Value())
	}

	var collections []string
	for k := range collectionspaths {
		collections = append(collections, string(collectionspaths[k]))
	}

	sort.Strings(collections)

	var clientCollections []string
	for k := range client.Collections {
		clientCollections = append(clientCollections, k)
	}

	sort.Strings(clientCollections)

	// WON'T FIX: This is OK. Not all collections on dbus created by a single client
	// if !reflect.DeepEqual(clientCollections, collections) {
	// 	panic(fmt.Sprintf("Client 'Collections' property is out of sync. Object: %v, dbus: %v",
	// 		clientCollections, collections))
	// }

	return collections, nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Properties <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
