package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"github.com/yousefvand/secret-service/pkg/crypto"
)

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Entities >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

type Database struct {
	// Database version (used for backward compatibility)
	Version string `json:"version"`
	// TRUE if database is encrypted otherwise false
	Encrypted bool `json:"encrypted"`
	// All collections in this database
	Collections []DbCollection `json:"collections"`
}

// Collection's Parent is Database (root)
type DbCollection struct {
	// Collection object path on dbus
	ObjectPath dbus.ObjectPath `json:"objectPath"`
	// All items in this collection
	Items []DbItem `json:"items"`
	// Collection properties
	Properties map[string]string `json:"properties"`
	// Collection Alias
	Alias string `json:"alias"`
	// Collection Label
	Label string `json:"label"`
	// Is collection locked?
	Locked bool `json:"locked"`
	// Collection creation time (epoch)
	Created uint64 `json:"created"`
	// Collection modification time (epoch)
	Modified uint64 `json:"modified"`
	// RawProperties map[string]string `json:"rawProperties"`
	// DbusProperties prop.Properties         `json:"dbusProperties"`
}

// Item's Parent is Collection
type DbItem struct {
	// Item parent (collection) object path
	Parent dbus.ObjectPath `json:"parent"`
	// Item object path on dbus
	ObjectPath dbus.ObjectPath `json:"objectPath"`
	// Item properties
	Properties map[string]string `json:"properties"`
	// Item secret (wrapper around SecretApi)
	Secret DbSecret `json:"secret"`
	// Item lookup attributes
	LookupAttributes map[string]string `json:"lookupAttributes"`
	// Item label
	Label string `json:"label"`
	// Is item locked?
	Locked bool `json:"locked"`
	// Item creation time (epoch)
	Created uint64 `json:"created"`
	// Item modification time (epoch)
	Modified uint64 `json:"modified"`
	// RawProperties    map[string]string `json:"rawProperties"`
	// DbusProperties prop.Properties         `json:"dbusProperties"`
}

// Secret's Parent is Item
type DbSecret struct {
	// Secret parent (item)
	Parent dbus.ObjectPath `json:"parent"`
	// Secret without encryption
	SecretText string `json:"secretText"`
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Entities <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> RestoreData >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// RestoreData reads database and restores dbus objects
func RestoreData(service *Service) {

	dbFile := filepath.Join(service.Config.Home, "db.json")

	db := Unmarshal(dbFile)

	if db == nil { // no db file
		close(service.DbLoadedChan)
		return // fresh run
	}

	encrypted := db.Encrypted // database is encrypted
	masterPassword := os.Getenv("MASTERPASSWORD")

	if len(masterPassword) > 0 && len(masterPassword) != 32 {
		log.Warnf("MASTERPASWORD length is wrong. Expected 32, got %d", len(masterPassword))
		masterPassword = ""
	}

	if encrypted && len(masterPassword) != 32 {
		log.Panicf("Database is encrypted but cannot find a 32 character MASTERPASSWORD")
	}

	// Iterating db Collections
	for _, collectionValue := range db.Collections {

		var collection *Collection

		// ignore creating default collection
		if collectionValue.Alias == "default" {
			collection = service.GetCollectionByAlias("default")
		} else {
			collection = NewCollection(service)
			// collection.SetProperties(properties)
			collection.Alias = collectionValue.Alias
			collection.ObjectPath = collectionValue.ObjectPath
			collection.Label = collectionValue.Label

			// Set Collection properties
			for k, v := range collectionValue.Properties {
				collection.Properties[k] = dbus.MakeVariant(v)
			}

			collection.Locked = collectionValue.Locked
			collection.Created = collectionValue.Created
			collection.Modified = collectionValue.Modified
			service.AddCollection(collection, collection.Locked,
				collection.Created, collection.Modified, false)
			service.UpdatePropertyCollections()
		}

		// Collection Items
		for _, ItemValue := range collectionValue.Items {
			item := NewItem(collection)
			item.ObjectPath = ItemValue.ObjectPath

			// Set Item properties
			for k, v := range ItemValue.Properties {
				item.Properties[k] = dbus.MakeVariant(v)
			}

			// Set Item LookupAttributes
			for k, v := range ItemValue.LookupAttributes {
				item.LookupAttributes[k] = v
			}

			item.Locked = ItemValue.Locked
			item.Created = ItemValue.Created
			item.Modified = ItemValue.Modified

			item.Secret.SecretApi.ContentType = "text/plain"

			if encrypted {
				decrypted, err := crypto.DecryptAESCBC256(masterPassword, ItemValue.Secret.SecretText)
				if err != nil {
					if os.Getenv("ENV") != "TEST" {
						log.Panicf("Cannot decrypt database. Error: %v", err)
					}
				}
				item.Secret.PlainSecret = decrypted
			} else {
				item.Secret.PlainSecret = ItemValue.Secret.SecretText
			}

			collection.AddItem(item, false, false, item.Locked, item.Created, item.Modified, true)
			collection.UpdatePropertyCollectionItems()
		}

	}

	close(service.DbLoadedChan) // Signal database has loaded
	log.Info("Loading data finished successfully")
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< RestoreData <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> PersistData >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// PersistData makes dbus objects persistent to db as soon as they change
func PersistData(ctx context.Context, service *Service) {

	dbFile := filepath.Join(service.Config.Home, "db.json")
	dbLock := new(sync.Mutex)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			<-service.SaveSignalChan // blocking until signal receives

			dbLock.Lock()

			log.Infof("Saving database at: '%s'", dbFile)
			Marshal(service, dbFile)

			dbLock.Unlock()

		}
	}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< PersistData <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Marshal >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// Marshal converts dbus objects to JSON
func Marshal(service *Service, dbFile string) {

	encrypt := service.Config.EncryptDatabase
	masterPassword := os.Getenv("MASTERPASSWORD")

	if len(masterPassword) > 0 && len(masterPassword) != 32 {
		log.Warnf("MASTERPASWORD length is wrong. Expected 32, got %d", len(masterPassword))
		masterPassword = ""
	}

	if encrypt && len(masterPassword) != 32 {
		log.Panicf("Cannot encrypt database with a non 32 character MASTERPASSWORD")
	}

	db := Database{}
	db.Version = "0.1.0"
	db.Encrypted = encrypt
	db.Collections = []DbCollection{}

	service.CollectionsMutex.RLock()

	for _, collectionValue := range service.Collections {

		collectionValue.DataMutex.RLock()

		collection := DbCollection{}
		collection.Items = []DbItem{}
		collection.Properties = make(map[string]string)

		collection.ObjectPath = collectionValue.ObjectPath

		collectionValue.ItemsMutex.RLock()
		for _, itemValue := range collectionValue.Items {

			itemValue.DataMutex.RLock()

			item := DbItem{}
			item.Properties = make(map[string]string)

			item.Parent = collectionValue.ObjectPath
			item.ObjectPath = itemValue.ObjectPath

			itemValue.PropertiesMutex.RLock()
			for k, v := range itemValue.Properties {
				if val, ok := v.Value().(string); ok {
					item.Properties[k] = val
				}
			}
			itemValue.PropertiesMutex.RUnlock()

			itemValue.Secret.DataMutex.RLock()

			secret := DbSecret{}
			secret.Parent = itemValue.ObjectPath

			if encrypt {
				encrypted, err := crypto.EncryptAESCBC256(masterPassword, itemValue.Secret.PlainSecret)

				if err != nil {
					log.Panicf("Database encryption failed. Error: %v", err)
				}

				secret.SecretText = encrypted
			} else {
				secret.SecretText = itemValue.Secret.PlainSecret
			}

			item.Secret = secret

			itemValue.Secret.DataMutex.RUnlock()

			item.LookupAttributes = itemValue.LookupAttributes
			item.Label = itemValue.Label
			itemValue.LockMutex.Lock()
			item.Locked = itemValue.Locked
			itemValue.LockMutex.Unlock()
			item.Created = itemValue.Created
			item.Modified = itemValue.Modified

			collection.Items = append(collection.Items, item)

			itemValue.DataMutex.RUnlock()
		}
		collectionValue.ItemsMutex.RUnlock()

		for k, v := range collectionValue.Properties {
			if val, ok := v.Value().(string); ok {
				collection.Properties[k] = val
			}
		}

		collection.Alias = collectionValue.Alias
		collection.Label = collectionValue.Label
		collectionValue.LockMutex.Lock()
		collection.Locked = collectionValue.Locked
		collectionValue.LockMutex.Unlock()
		collection.Created = collectionValue.Created
		collection.Modified = collectionValue.Modified

		db.Collections = append(db.Collections, collection)

		collectionValue.DataMutex.RUnlock()
	}

	service.CollectionsMutex.RUnlock()

	// save db to file

	content, err := json.MarshalIndent(db, "", " ")
	if err != nil {
		log.Panicf("Cannot marshal database. Error: %v", err)
	}

	err = os.WriteFile(dbFile, content, 0600)
	if err != nil {
		log.Panicf("Cannot write to database. Error: %v", err)
	}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Marshal <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Unmarshal >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// Unmarshal converts JSON data into dbus objects
func Unmarshal(dbFile string) *Database {

	dbExist, err := fileOrFolderExists(dbFile)

	if err != nil {
		log.Panicf("Cannot check db file existence at: '%s'. Error: %v", dbFile, err)
	}

	// This is a fresh run, no db exist yet
	if dbExist {
		log.Infof("Loading data from: '%s'", dbFile)
	} else {
		return nil
	}

	content, err := os.ReadFile(dbFile)

	if err != nil {
		log.Panicf("Cannot read database file at '%s'. Error: %v", dbFile, err)
	}

	var db Database

	err = json.Unmarshal(content, &db)

	if err != nil {
		log.Panicf("Malformed database at '%s'. Error: %v", dbFile, err)
	}

	return &db

}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Unmarshal <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

// fileOrFolderExists returns true if a file/folder exist otherwise false
func fileOrFolderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
