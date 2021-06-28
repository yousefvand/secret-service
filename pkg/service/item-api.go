// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

/////////////////////////////////// Methods ///////////////////////////////////

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Delete >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Delete ( OUT ObjectPath Prompt);
*/

// Delete removes an item from a collection
func (item *Item) Delete() (dbus.ObjectPath, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "org.freedesktop.Secret.Item",
		"method":    "Delete",
		"item path": item.ObjectPath,
	}).Trace("Method called by client")

	item.Parent.RemoveItem(item)

	item.SignalItemDeleted()
	item.Parent.UpdatePropertyCollectionItems()
	item.Parent.UpdateModified()

	// A prompt object, or the special value ‘/’ if no prompt is necessary.
	return dbus.ObjectPath("/"), nil

}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Delete <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetSecret >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	GetSecret ( IN ObjectPath session,
	            OUT Secret secret);
*/

// GetSecret retrieves the secret for this item
func (item *Item) GetSecret(session dbus.ObjectPath) (*SecretApi, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "org.freedesktop.Secret.Item",
		"method":    "GetSecret",
		"session":   session,
	}).Trace("Method called by client")

	secretApi := &SecretApi{}
	service := item.Parent.Parent
	sessionInUse := service.GetSessionByPath(session)

	if sessionInUse == nil {
		log.Warn("Secret session is missing")
		// TODO: Check if dbus get corrupted by returning dbus error
		return nil, ApiErrorNoSession() // empty secretApi
	}

	secretApi.Session = session
	if sessionInUse.EncryptionAlgorithm == Plain {
		secretApi.Value = []byte(item.Secret.PlainSecret)
		secretApi.ContentType = "text/plain; charset=utf8"
		secretApi.Parameters = []byte("")
	} else { // dh-ietf1024-sha256-aes128-cbc-pkcs7
		iv, cipherData, err := AesCBCEncrypt([]byte(item.Secret.PlainSecret),
			sessionInUse.SymmetricKey)
		if err != nil {
			log.Errorf("Cannot GetSecret due to encryption error. Error: %v", err)
			return nil, // empty secret
				DbusErrorCallFailed("Cannot GetSecret due to encryption error. Error: " + err.Error())
		}
		secretApi.Parameters = iv
		secretApi.Value = cipherData
	}
	log.Tracef("GetSecret returned: %s", item.Secret.PlainSecret)

	return secretApi, nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< GetSecret <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> SetSecret >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	SetSecret ( IN Secret secret);
*/

// SetSecret sets the secret for this item
func (item *Item) SetSecret(secretApi SecretApi) *dbus.Error {

	log.WithFields(log.Fields{
		"interface": "org.freedesktop.Secret.Item",
		"method":    "SetSecret",
		"secretApi": secretApi,
	}).Trace("Method called by client")

	secret := NewSecret(item)
	session := item.Parent.Parent.GetSessionByPath(secretApi.Session)

	if session == nil {
		log.Warn("Secret session is missing")
		return ApiErrorNoSession() // TODO: Check if dbus get corrupted by returning dbus error
	}

	secret.SecretApi = &secretApi

	if session.EncryptionAlgorithm == Plain {
		secret.PlainSecret = string(secret.SecretApi.Value)
	} else {
		iv := secret.SecretApi.Parameters
		plainSecret, err := AesCBCDecrypt(iv, secret.SecretApi.Value, session.SymmetricKey)
		if err != nil {
			log.Errorf("Cannot SetSecret due to decryption error. Error: %v", err)
			return DbusErrorCallFailed("Cannot SetSecret due to decryption error. Error: " + err.Error())
		}
		secret.PlainSecret = string(plainSecret)
	}

	item.DataMutex.Lock()
	item.Secret = secret
	item.DataMutex.Unlock()
	// TODO: Update dbus?
	item.SignalItemChanged()
	item.Parent.UpdateModified()
	log.Tracef("SetSecret received: %s", item.Secret.PlainSecret)
	item.SaveData()

	return nil
}
