// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// create and initialize a new instance of service
func New() *Service {
	service := new(Service)
	service.Connection = nil
	// service.LockMutex = new(sync.Mutex) // TODO: Remove
	service.SessionsMutex = new(sync.RWMutex)
	service.CollectionsMutex = new(sync.RWMutex)
	service.Sessions = make(map[string]*Session)
	service.DbLoadedChan = make(chan struct{})
	service.SaveSignalChan = make(chan struct{})
	service.Collections = make(map[string]*Collection)
	service.ServiceReadyChan = make(chan struct{})
	service.ServiceShutdownChan = make(chan struct{})
	service.SecretService = &SecretService{}
	service.SecretService.Session = &SecretServiceCLiSession{}
	service.SecretService.Parent = service
	// service.SecretService.Session = &SecretServiceCLiSession{}

	// service.Update callback is set by the user (App)
	// Service.SaveData = func() {
	// 	Service.SaveSignalChan <- struct{}{}
	// }
	return service
}

// send signal to save database
func (s *Service) SaveData() {
	s.SaveSignalChan <- struct{}{}
}

// connect to session dbus
func (s *Service) connect() {
	connection, err := dbus.SessionBus()
	if err != nil {
		log.Panicf("Cannot connect to session dbus. Error: %v", err)
	}
	s.Connection = connection
}

// disconnect drops connection to session dbus
func (s *Service) disconnect() {
	err := s.Connection.Close()
	if err != nil {
		log.Panic("Cannot disconnect from dbus")
	}
}

// start secretserviced
func (service *Service) Start(ctx context.Context) {

	log.Info("===== Secret Service Started =====")
	log.Info("Secret service dbus address: /org/freedesktop/secrets")
	log.Debugf("Using total of %v MiB of OS memory", MemUsageOS())
	log.Debugf("Service is using %d Goroutines", runtime.NumGoroutine())

	if service.Connection == nil {
		service.connect()
	}

	getRootName(service.Connection)    // own 'org.freedesktop.secrets' on dbus
	dbusInitialize(service.Connection) // make initial dbus objects (i.e. /org)
	/* create deafult collection at: '/org/freedesktop/secrets/aliases/default' */
	epoch := Epoch()
	DefaultCollection(service, false, epoch, epoch)
	dbusSecretService(service) // TODO: temp
	// create SecretService interface on dbus path: '/org/freedesktop/secrets'
	dbusService(service)

	go RestoreData(service)
	<-service.DbLoadedChan
	go PersistData(ctx, service)

	close(service.ServiceReadyChan) // propagate a signal that means service is ready

	<-ctx.Done() // waiting for shutdown signal
	service.disconnect()
	log.Info("===== Secret Service gracefully shutted down =====")
	close(service.ServiceShutdownChan)
}

// locks the service
// func (s *Service) ServiceLock() {
// 	s.LockMutex.Lock()
// 	s.Locked = true
// 	s.LockMutex.Unlock()
// }

// unlocks the service
// func (s *Service) ServiceUnlock() {
// 	s.LockMutex.Lock()
// 	s.Locked = false
// 	s.LockMutex.Unlock()
// }

// add a new session to service's session map
func (s *Service) AddSession(session *Session) {
	s.SessionsMutex.Lock()
	s.Sessions[string(session.ObjectPath)] = session
	s.SessionsMutex.Unlock()
	// add session object to dbus
	dbusAddSession(s, session)
	log.Infof("New session at: %v", session.ObjectPath)
	// s.SaveData()
}

// remove a session from service's session map
func (s *Service) RemoveSession(session *Session) {
	s.SessionsMutex.Lock()
	_, ok := s.Sessions[string(session.ObjectPath)]
	if !ok {
		log.Errorf("Session doesn't exist to be removed: %v",
			session.ObjectPath)
		return
	}
	delete(s.Sessions, string(session.ObjectPath))
	s.SessionsMutex.Unlock()
	// update dbus objects after session is removed
	dbusUpdateSessions(s)
	s.SaveData()
	log.Infof("Session removed: %v", session.ObjectPath)
}

// HasSession returns true if session exists otherwise false
func (s *Service) HasSession(sessionPath dbus.ObjectPath) bool {
	s.SessionsMutex.RLock()
	_, ok := s.Sessions[string(sessionPath)]
	s.SessionsMutex.RUnlock()
	return ok
}

// GetSessionByPath returns session with given objectpath
func (service *Service) GetSessionByPath(sessionPath dbus.ObjectPath) *Session {
	service.SessionsMutex.RLock()
	defer service.SessionsMutex.RUnlock()

	for _, session := range service.Sessions {
		if session.ObjectPath == sessionPath {
			return session
		}
	}
	return nil
}

// add a new collection to service's collection map
func (s *Service) AddCollection(collection *Collection,
	locked bool, created uint64, modified uint64, saveData bool) {

	s.CollectionsMutex.Lock()
	s.Collections[string(collection.ObjectPath)] = collection
	s.CollectionsMutex.Unlock()
	dbusAddCollection(collection, locked, created, modified)

	// Let database be loaded before saving anything
	if collection.Alias != "default" && saveData {
		s.SaveData()
	}
}

// remove a collection from service's collection map
func (s *Service) RemoveCollection(collection *Collection) {
	s.CollectionsMutex.Lock()
	_, ok := s.Collections[string(collection.ObjectPath)]
	if !ok {
		log.Errorf("Collection doesn't exist to be removed: %v",
			collection.ObjectPath)
		return
	}
	delete(s.Collections, string(collection.ObjectPath))
	s.CollectionsMutex.Unlock()
	dbusUpdateCollections(s)
	log.Infof("Collection removed: %v", collection.ObjectPath)
	s.SaveData()
}

// HasCollection returns true if collection exists otherwise false
func (s *Service) HasCollection(collectionPath dbus.ObjectPath) bool {
	s.CollectionsMutex.RLock()
	_, ok := s.Collections[string(collectionPath)]
	s.CollectionsMutex.RUnlock()
	return ok
}

func (service *Service) GetCollectionByPath(collectionPath dbus.ObjectPath) *Collection {
	for _, collection := range service.Collections {
		if collection.ObjectPath == collectionPath {
			return collection
		}
	}
	return nil
}

// GetCollectionByAlias finds and return a collection by it's alias name otherwise return nil
func (s *Service) GetCollectionByAlias(alias string) *Collection {

	if alias == "" {
		return nil
	}

	result := []*Collection{}

	s.CollectionsMutex.RLock()
	for _, v := range s.Collections {
		if v.Alias == alias {
			result = append(result, v)
		}
	}
	s.CollectionsMutex.RUnlock()

	if len(result) > 1 {
		log.Panicf("There are %d collections with the same alias '%s'", len(result), alias)
	} else if len(result) == 1 {
		return result[0]
	}
	return nil

}

func (service *Service) GetItemByPath(itemPath dbus.ObjectPath) *Item {
	for _, collection := range service.Collections {
		for _, item := range collection.Items {
			if item.ObjectPath == itemPath {
				return item
			}
		}
	}
	return nil
}

// UpdatePropertyCollections updates dbus properties of Service
func (s *Service) UpdatePropertyCollections() {

	var collections []string

	s.CollectionsMutex.RLock()
	defer s.CollectionsMutex.RUnlock()

	for _, collection := range s.Collections {
		collections = append(collections, string(collection.ObjectPath))
	}

	PropsService.SetMust("org.freedesktop.Secret.Service",
		"Collections", collections)
}

// ReadPasswordFile returns contents of 'password.yaml' file if exists otherwise empty string
func (service *Service) ReadPasswordFile() string {

	passwordFilePath := filepath.Join(service.Home, "password.yaml")
	exist, err := fileOrFolderExists(passwordFilePath)

	if err != nil {
		log.Panicf("Cannot determine 'password.yaml' file status: %s", passwordFilePath)
	}

	if !exist {
		return ""
	}

	data, err := ioutil.ReadFile(passwordFilePath)

	if err != nil {
		log.Warnf("Cannot open 'password.yaml' file at: %s.", passwordFilePath)
	}

	var passwordFile PasswordFile
	err = yaml.Unmarshal(data, &passwordFile)

	if err != nil {
		log.Warnf("found malformed 'password.yaml' file: '%s'.", passwordFilePath)
	}

	if len(passwordFile.PasswordHash) < 1 {
		return ""
	}

	return passwordFile.PasswordHash

}

// WritePasswordFile writes 'password.yaml' file or returns error
func (service *Service) WritePasswordFile(passwordHash string) error {

	const version string = "0.1.0"

	var content []byte = []byte(`# Password file version
version: ` + version + `
# Password hash: sha512(salt+password)
passwordHash: '` + passwordHash + `'`)

	passwordFile := filepath.Join(service.Home, "password.yaml")
	errWritePasswordFile := ioutil.WriteFile(passwordFile, content, 0600)

	if errWritePasswordFile != nil {
		log.Warnf("Cannot write password file. Error: %v", errWritePasswordFile)
		return errWritePasswordFile
	}

	return nil
}
