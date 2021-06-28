// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/crypto/hkdf"

	"github.com/godbus/dbus/v5"
	"github.com/monnand/dhkx"
	log "github.com/sirupsen/logrus"
)

/*
	API implementation of:
	org.freedesktop.Secret.Service
*/

/////////////////////////////////// Methods ///////////////////////////////////

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OpenSession >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	OpenSession ( IN String algorithm,
	              IN Variant input,
	              OUT Variant output,
	              OUT ObjectPath result);
*/

// OpenSession opens a unique session for the caller application
// further communication encryption/decryption relies on the related session
func (service *Service) OpenSession(algorithm string,
	input dbus.Variant) (dbus.Variant, dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "org.freedesktop.Secret.Service",
		"method":    "OpenSession",
		"algorithm": algorithm,
		"input":     input.Value(),
	}).Trace("Method called by client")

	// TODO: Remove
	// if service.Locked {
	// 	log.Warn("Cannot 'OpenSession' when service is locked.")
	// 	return dbus.MakeVariant(""), dbus.ObjectPath("/"), ApiErrorIsLocked()
	// }

	log.Debugf("Client suggested '%s' algorithm", algorithm)

	// if OpenSession succeeds all related information are stored in a session
	session := NewSession(service)

	switch strings.ToLower(algorithm) {

	case "plain":
		session.EncryptionAlgorithm = Plain

		switch t := input.Value().(type) {

		case string:
			if (input.Value()).(string) != "" { // we expect an empty string with plain algorithm
				log.Warnf("Wrong call parameter to 'OpenSession' method. While 'plain' algorithm selected, 'input' parameter is not empty: %s", (input.Value()).(string))
				return dbus.MakeVariant(""), dbus.ObjectPath("/"),
					DbusErrorInvalidArgs("expected empty string for 'input', got: " + (input.Value()).(string))
			}

		default: // input type is not string. This is a violation of API
			log.Errorf("Wrong call parameter to 'OpenSession' method. Unknown 'input' type: %T", t)
			return dbus.MakeVariant(""), dbus.ObjectPath("/"),
				DbusErrorInvalidSignature(fmt.Sprintf("expected 'string' for 'input' got: %T, with value: %v", input.Value(), input.Value()))

		}

		log.Debug("Agreed on 'plain' algorithm")
		path := "/org/freedesktop/secrets/session/" + UUID()
		session.ObjectPath = dbus.ObjectPath(path)

		service.AddSession(session)

		return dbus.MakeVariant(""), session.ObjectPath, nil // end of successful negotiation

	case "dh-ietf1024-sha256-aes128-cbc-pkcs7":
		session.EncryptionAlgorithm = Dh_ietf1024_sha256_aes128_cbc_pkcs7

		group, err := dhkx.GetGroup(2) // 2 -> 1024 bit (128 bytes) secret key
		if err != nil {
			log.Panicf("Diffie–Hellman group creation failed. Error: %s", err.Error())
			return dbus.MakeVariant(""), dbus.ObjectPath("/"),
				DbusErrorCallFailed("Diffie–Hellman group creation failed. Error: " + err.Error())
		}
		session.Group = group // TODO: Remove me (not needed)
		privateKey, err := group.GeneratePrivateKey(rand.Reader)
		if err != nil {
			log.Panicf("Diffie–Hellman private key generation failed. Error: %s", err.Error())
			return dbus.MakeVariant(""), dbus.ObjectPath("/"),
				DbusErrorCallFailed("Diffie–Hellman private key generation failed. Error: " + err.Error())
		}
		session.PrivateKey = privateKey // // TODO: Remove me (not needed)

		// TODO: big endian
		/*
			b := []byte{...}
			for i := 0; i < len(b)/2; i++ {
			    b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
			}
		*/
		session.PublicKey = privateKey.Bytes() // TODO: Remove me (not needed)

		// TODO: Check convertion validity
		var clientPublicKey []byte
		err = dbus.Store([]interface{}{input.Value()}, &clientPublicKey)
		if err != nil {
			log.Panicf("Cannot convert client public key. Error: %v", err)
			return dbus.MakeVariant(""), dbus.ObjectPath("/"),
				DbusErrorCallFailed("Cannot convert client public key. Error: " + err.Error())
		}

		log.Tracef("Client public key: %v", clientPublicKey)
		log.Tracef("Client public key length: %v", len(clientPublicKey))

		if len(clientPublicKey) != 128 {
			log.Errorf("Invalid client public key length. Expected 128 bytes, received %v bytes.",
				len(clientPublicKey))
			return dbus.MakeVariant(""), dbus.ObjectPath("/"),
				DbusErrorInvalidArgs("Wrong length of public key. Expected 128 bytes got " +
					strconv.Itoa(len(clientPublicKey)) + " bytes")
		}

		// TODO: Remove me (not needed)
		session.ClientPublicKey = clientPublicKey // dhkx.NewPublicKey(clientPublicKey)

		sharedKey, err := group.ComputeKey(dhkx.NewPublicKey(clientPublicKey), session.PrivateKey)
		if err != nil {
			log.Panicf("Diffie–Hellman shared key generation failed. Error: %s", err.Error())
			return dbus.MakeVariant(""), dbus.ObjectPath("/"),
				DbusErrorCallFailed("Diffie–Hellman shared key generation failed. Error: " + err.Error())
		}
		// TODO: Remove me (not needed)
		session.SharedKey = sharedKey.Bytes()

		log.Tracef("Shared key: %v", session.SharedKey)
		log.Tracef("Shared key length: %v", len(session.SharedKey))

		hkdf := hkdf.New(sha256.New, session.SharedKey, nil, nil)
		symmetricKey := make([]byte, aes.BlockSize) // 16 * 8 = 128 bit
		n, err := io.ReadFull(hkdf, symmetricKey)
		if n != aes.BlockSize {
			log.Panicf("Cannot create 16 byte key. Length is: %v", len(symmetricKey))
			return dbus.MakeVariant(""), dbus.ObjectPath("/"),
				DbusErrorCallFailed("Cannot create 16 byte key. Length is: " + strconv.Itoa(len(symmetricKey)))
		}
		if err != nil {
			log.Panicf("Symmetric Key generation failed. Error: %s", err.Error())
			return dbus.MakeVariant(""), dbus.ObjectPath("/"),
				DbusErrorCallFailed("Symmetric Key generation failed. Error: " + err.Error())
		}
		session.SymmetricKey = symmetricKey

		log.Tracef("Symmetric key: %v", session.SymmetricKey)
		log.Tracef("Symmetric key length: %v", len(session.SymmetricKey))

		log.Debug("Agreed on 'dh-ietf1024-sha256-aes128-cbc-pkcs7' algorithm")

		path := "/org/freedesktop/secrets/session/" + UUID()
		session.ObjectPath = dbus.ObjectPath(path)

		service.AddSession(session)

		// TODO: make sure it is big endian
		return dbus.MakeVariant(session.PublicKey), dbus.ObjectPath(path), nil // end of successful negotiation

	default: // algorithm is not 'plain' or 'dh-ietf1024-sha256-aes128-cbc-pkcs7'
		log.Warnf("The '%s' algorithm suggested by client is not supported", algorithm)
		return dbus.MakeVariant(""), dbus.ObjectPath("/"), ApiErrorNotSupported()

	}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< OpenSession <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CreateCollection >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	CreateCollection ( IN Dict<String,Variant> properties,
	                   IN String alias,
	                   OUT ObjectPath collection,
	                   OUT ObjectPath prompt);
*/

// CreateCollection creates a collection which can hold multiple items
func (service *Service) CreateCollection(properties map[string]dbus.Variant,
	alias string) (dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":  "org.freedesktop.Secret.Service",
		"method":     "CreateCollection",
		"properties": properties,
		"alias":      alias,
	}).Trace("Method called by client")

	if len(properties) == 0 { // FIXME: return error
		log.Warn("Client asked to create a collection with empty 'properties'")
	}

	// if a collection with the same alias exist return that
	collection := service.GetCollectionByAlias(alias)

	// no collection with the same alias. Let's create one:
	if collection == nil {
		collection = NewCollection(service)
		collection.SetProperties(properties)
		collection.Alias = strings.TrimSpace(alias)

		uuid := UUID()
		collectionLabel := uuid[len(uuid)/2:] // use the last half of UUID
		// Use org.freedesktop.Secret.Collection.Label if available
		if collection.Label != "" {
			// Make sure path doesn't exist
			if service.GetCollectionByPath(
				dbus.ObjectPath("/org/freedesktop/secrets/collection/"+collection.Label),
			) == nil {
				collectionLabel = collection.Label // use label in collection path (override uuid)
			}
		}
		// path should be '/' if prompting is necessary
		path := "/org/freedesktop/secrets/collection/" + collectionLabel
		collection.ObjectPath = dbus.ObjectPath(path)

		epoch := Epoch()
		service.AddCollection(collection, false, epoch, epoch, true)
		collection.SignalCollectionCreated() // TODO: if it is default collection -> no signal

		if collection.Alias == "" {
			log.Infof("New collection with no alias at: %v", collection.ObjectPath)
		} else {
			log.Infof("New collection with alias '%s' at: %v", collection.Alias, collection.ObjectPath)
		}
	} else if collection.Alias != "default" {
		log.Infof("Collection with alias '%s' already exists at: %v", collection.Alias, collection.ObjectPath)
	}

	if collection.Alias == "default" {
		log.Info("Client asked for default collection at: /org/freedesktop/secrets/aliases/default")
	}

	service.UpdatePropertyCollections()

	// A prompt object if prompting is necessary, or ‘/’ if no prompt was needed.
	return collection.ObjectPath, dbus.ObjectPath("/"), nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< CreateCollection <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> SearchItems >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	SearchItems ( IN Dict<String,String> attributes,
	              OUT Array<ObjectPath> unlocked,
	              OUT Array<ObjectPath> locked);
*/

// SearchItems finds items inside all collection. A collection
// consists of many items: item = secret + lookup attributes + label
func (service *Service) SearchItems(
	attributes map[string]string) ([]dbus.ObjectPath, []dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":  "org.freedesktop.Secret.Service",
		"method":     "SearchItems",
		"attributes": attributes,
	}).Trace("Method called by client")

	var lockedItems []dbus.ObjectPath
	var unlockedItems []dbus.ObjectPath

	for _, collection := range service.Collections {
		for _, item := range collection.Items {
			// FIXME: Singe or Full match (FullMatch works with skype)
			if IsMapSubsetFullMatch(item.LookupAttributes, // FIXME SingleMatch or FullMatch
				attributes, collection.ItemsMutex) {

				log.Debugf("SearchItems found match. Label: %s, Path: %s", item.Label, item.ObjectPath)
				if item.Locked {
					lockedItems = append(lockedItems, item.ObjectPath)
				} else {
					unlockedItems = append(unlockedItems, item.ObjectPath)
				}
			} else {
				log.Debugf("SearchItems didn't find any match for: %v", attributes)
			}
		}
	}

	log.WithFields(log.Fields{
		"unlockedItems": unlockedItems,
		"lockedItems":   lockedItems,
	}).Tracef("SearchItems results for: %v", attributes)

	return unlockedItems, lockedItems, nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< SearchItems <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Unlock >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Unlock ( IN Array<ObjectPath> objects,
	         OUT Array<ObjectPath> unlocked,
	         OUT ObjectPath prompt);
*/

// Unlock unlocks the specified objects (collections, items)
func (service *Service) Unlock(
	objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "org.freedesktop.Secret.Service",
		"method":    "Unlock",
		"objects":   objects,
	}).Trace("Method called by client")

	var unlockedObjects []dbus.ObjectPath

	for _, object := range objects {
		for _, collection := range service.Collections {
			if collection.ObjectPath == object {
				if collection.Locked {
					collection.Unlock()
					collection.UpdateModified()
					collection.SignalCollectionChanged()
					unlockedObjects = append(unlockedObjects, collection.ObjectPath)
				}
			}
			for _, item := range collection.Items {
				if item.ObjectPath == object {
					if item.Locked {
						item.Unlock()
						item.UpdateModified()
						item.SignalItemChanged()
						unlockedObjects = append(unlockedObjects, item.ObjectPath)
					}
				}
			}
		}
	}
	log.Debugf("Unlocked objects: %v", unlockedObjects)
	service.SaveData()

	return unlockedObjects, dbus.ObjectPath("/"), nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Unlock <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Lock >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Lock ( IN Array<ObjectPath> objects,
	       OUT Array<ObjectPath> locked,
	       OUT ObjectPath Prompt);
*/

// Lock locks the specified objects (collections, items)
func (service *Service) Lock(
	objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "org.freedesktop.Secret.Service",
		"method":    "Lock",
		"objects":   objects,
	}).Trace("Method called by client")

	var lockedObjects []dbus.ObjectPath

	for _, object := range objects {
		for _, collection := range service.Collections {
			if collection.ObjectPath == object {
				if !collection.Locked {
					collection.Lock()
					// FIXME: EMit signal 'SignalCollectionChanged' and update modified time?
					lockedObjects = append(lockedObjects, collection.ObjectPath)
				}
			}
			for _, item := range collection.Items {
				if item.ObjectPath == object {
					if !item.Locked {
						item.Lock()
						// FIXME: EMit signal 'SignalItemChanged' and update modified time?
						lockedObjects = append(lockedObjects, item.ObjectPath)
					}
				}
			}
		}
	}
	log.Debugf("Locked objects: %v", lockedObjects)
	service.SaveData()

	return lockedObjects, dbus.ObjectPath("/"), nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Lock <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetSecrets >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	GetSecrets ( IN Array<ObjectPath> items,
	             IN ObjectPath session,
	             OUT Dict<ObjectPath,Secret> secrets);
*/

// GetSecrets retrieves multiple secrets from different items
func (service *Service) GetSecrets(items []dbus.ObjectPath,
	session dbus.ObjectPath) (map[dbus.ObjectPath]SecretApi, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "org.freedesktop.Secret.Service",
		"method":    "GetSecrets",
		"items":     items,
		"session":   session,
	}).Trace("Method called by client")

	sessionInUse := service.GetSessionByPath(session)

	if sessionInUse == nil {
		log.Warnf("GetSecrets session doesn't exist: %s", string(session))
		return map[dbus.ObjectPath]SecretApi{}, ApiErrorNoSession()
	}

	result := make(map[dbus.ObjectPath]SecretApi)

	for _, collection := range service.Collections {
		for _, item := range collection.Items {
			for _, itemPath := range items {
				if item.ObjectPath == itemPath {
					secretApi := SecretApi{}
					secretApi.Session = sessionInUse.ObjectPath
					secretApi.ContentType = item.Secret.SecretApi.ContentType
					iv, cipherData, err := AesCBCEncrypt([]byte(item.Secret.PlainSecret),
						[]byte(sessionInUse.SymmetricKey))
					if err != nil {
						log.Errorf("Cannot GetSecrets due to encryption error. Error: %v", err)
						return map[dbus.ObjectPath]SecretApi{},
							DbusErrorCallFailed("Cannot GetSecrets due to encryption error. Error: " + err.Error())
					}
					secretApi.Parameters = iv
					secretApi.Value = cipherData
					result[itemPath] = secretApi

					log.WithFields(log.Fields{ // TODO: Remove me.
						"item plain secret":   item.Secret.PlainSecret,
						"item secretApi":      item.Secret.SecretApi,
						"generated secretApi": secretApi,
					}).Debug("FIXME: Are they the same?")
				}
			}
		}
	}
	log.Tracef("GetSecrets result: %v", result)

	return result, nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< GetSecrets <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ReadAlias >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	ReadAlias ( IN String name,
	            OUT ObjectPath collection);
*/

// returns the collection with the given alias
func (service *Service) ReadAlias(name string) (dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":    "org.freedesktop.Secret.Service",
		"method":       "ReadAlias",
		"name (alias)": name,
	}).Trace("Method called by client")

	if collection := service.GetCollectionByAlias(name); collection != nil {
		log.Infof("Found alias '%v' for collection: %v", name, collection.ObjectPath)
		return collection.ObjectPath, nil
	}

	log.Infof("Found no collection with alias: %v", name)
	return dbus.ObjectPath("/"), nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< ReadAlias <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> SetAlias >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	SetAlias ( IN String name,
	           IN ObjectPath collection);
*/

// set an alias for a collection
func (service *Service) SetAlias(name string, collection dbus.ObjectPath) *dbus.Error {

	log.WithFields(log.Fields{
		"interface":  "org.freedesktop.Secret.Service",
		"method":     "SetAlias",
		"name":       name,
		"collection": collection,
	}).Trace("Method called by client")

	if collection == "/org/freedesktop/secrets/aliases/default" {
		log.Warnf("Client tried to change 'default' collection alias to '%v'", name)
		return ApiErrorNotSupported()
	}

	if c := service.GetCollectionByPath(collection); c != nil {
		if name == "/" {
			c.DataMutex.Lock()
			c.Alias = ""
			c.DataMutex.Unlock()
			log.Infof("Removed alias '%v' from collection: %v", c.Alias, c.ObjectPath)
		} else {
			c.DataMutex.Lock()
			c.Alias = name
			c.DataMutex.Unlock()
			log.Infof("Changed alias '%v' to '%v' for collection: %v", c.Alias, name, c.ObjectPath)
		}
		c.UpdateModified()
		c.SignalCollectionChanged()
		service.SaveData()

		return nil
	}

	return ApiErrorNoSuchObject()
}
