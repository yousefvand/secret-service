// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"strconv"
	"strings"

	"golang.org/x/crypto/hkdf"

	"github.com/godbus/dbus/v5"
	"github.com/monnand/dhkx"
	log "github.com/sirupsen/logrus"
)

/////////////////////////////////// Methods ///////////////////////////////////

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OpenSession >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	CreateSession ( IN String algorithm,
	                IN Variant input,
	                OUT Variant output,
	                OUT String serialnumber);
*/

// OpenSession opens a unique session for the caller application
// further communication encryption/decryption relies on the related session
func (service *Service) CreateSession(algorithm string,
	input dbus.Variant) (dbus.Variant, string, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "ir.remisa.SecretService",
		"method":    "CreateSession",
		"algorithm": algorithm,
		"input":     input.Value(),
	}).Trace("Method called by client")

	log.Debugf("Client suggested '%s' algorithm", algorithm)

	switch strings.ToLower(algorithm) {

	case "dh-ietf1024-sha256-aes128-cbc-pkcs7":

		group, err := dhkx.GetGroup(2) // 2 -> 1024 bit (128 bytes) secret key
		if err != nil {
			log.Panicf("Diffie–Hellman group creation failed. Error: %s", err.Error())
			return dbus.MakeVariant(""), "",
				DbusErrorCallFailed("Diffie–Hellman group creation failed. Error: " + err.Error())
		}

		privateKey, err := group.GeneratePrivateKey(rand.Reader)
		if err != nil {
			log.Panicf("Diffie–Hellman private key generation failed. Error: %s", err.Error())
			return dbus.MakeVariant(""), "",
				DbusErrorCallFailed("Diffie–Hellman private key generation failed. Error: " + err.Error())
		}

		publicKey := privateKey.Bytes()

		var clientPublicKey []byte
		err = dbus.Store([]interface{}{input.Value()}, &clientPublicKey)
		if err != nil {
			log.Panicf("Cannot convert client public key. Error: %v", err)
			return dbus.MakeVariant(""), "",
				DbusErrorCallFailed("Cannot convert client public key. Error: " + err.Error())
		}

		log.Tracef("Client public key: %v", clientPublicKey)
		log.Tracef("Client public key length: %v", len(clientPublicKey))

		if len(clientPublicKey) != 128 {
			log.Errorf("Invalid client public key length. Expected 128 bytes, received %v bytes.",
				len(clientPublicKey))
			return dbus.MakeVariant(""), "",
				DbusErrorInvalidArgs("Wrong length of public key. Expected 128 bytes got " +
					strconv.Itoa(len(clientPublicKey)) + " bytes")
		}

		sharedKey, err := group.ComputeKey(dhkx.NewPublicKey(clientPublicKey), privateKey)
		if err != nil {
			log.Panicf("Diffie–Hellman shared key generation failed. Error: %s", err.Error())
			return dbus.MakeVariant(""), "",
				DbusErrorCallFailed("Diffie–Hellman shared key generation failed. Error: " + err.Error())
		}

		sessionSharedKey := sharedKey.Bytes()

		log.Tracef("Shared key: %v", sessionSharedKey)
		log.Tracef("Shared key length: %v", len(sessionSharedKey))

		hkdf := hkdf.New(sha256.New, sessionSharedKey, nil, nil)
		symmetricKey := make([]byte, aes.BlockSize) // 16 * 8 = 128 bit
		n, err := io.ReadFull(hkdf, symmetricKey)
		if n != aes.BlockSize {
			log.Panicf("Cannot create 16 byte key. Length is: %v", len(symmetricKey))
			return dbus.MakeVariant(""), "",
				DbusErrorCallFailed("Cannot create 16 byte key. Length is: " + strconv.Itoa(len(symmetricKey)))
		}
		if err != nil {
			log.Panicf("Symmetric Key generation failed. Error: %s", err.Error())
			return dbus.MakeVariant(""), "",
				DbusErrorCallFailed("Symmetric Key generation failed. Error: " + err.Error())
		}

		service.SecretService.Session.SymmetricKey = symmetricKey

		log.Tracef("Symmetric key: %v", symmetricKey)
		log.Tracef("Symmetric key length: %v", len(symmetricKey))

		log.Debug("Agreed on 'dh-ietf1024-sha256-aes128-cbc-pkcs7' algorithm")

		service.SecretService.Session.SerialNumber = UUID()

		return dbus.MakeVariant(publicKey), service.SecretService.Session.SerialNumber, nil // end of successful negotiation

	default: // algorithm is not 'plain' or 'dh-ietf1024-sha256-aes128-cbc-pkcs7'
		log.Warnf("The '%s' algorithm suggested by client is not supported", algorithm)
		return dbus.MakeVariant(""), "", ApiErrorNotSupported()

	}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< OpenSession <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> SetPassword >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	SetPassword ( IN  String      serialnumber
		            IN  Array<Byte> oldpassword,
								IN  Array<Byte> oldpassword_iv,
								IN  Array<Byte> newpassword,
								IN  Array<Byte> newpassword_iv,
								IN  Array<Byte> oldSalt,
								IN  Array<Byte> oldSalt_iv,
								IN  Array<Byte> newSalt,
								IN  Array<Byte> newSalt_iv
								OUT String result);
*/

// Set password for first time or change a password of service
func (service *Service) SetPassword(serialnumber string,
	oldPassword []byte, oldPassword_iv []byte,
	newPassword []byte, newPassword_iv []byte,
	oldSalt []byte, oldSalt_iv []byte,
	newSalt []byte, newSalt_iv []byte) (string, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":      "ir.remisa.SecretService",
		"method":         "SetPassword",
		"serialnumber":   serialnumber,
		"oldPassword":    oldPassword,
		"oldPassword_iv": oldPassword_iv,
		"newPassword":    newPassword,
		"newPassword_iv": newPassword_iv,
		"oldSalt":        oldSalt,
		"oldSalt_iv":     oldSalt_iv,
		"newSalt":        newSalt,
		"newSalt_iv":     newSalt_iv,
	}).Trace("Method called by client")

	if service.SecretService.Session.SerialNumber != serialnumber {
		log.Warnf("Session mismatch: Expected: %s, got: %s",
			service.SecretService.Session.SerialNumber, serialnumber)
		return "session mismatch", nil
	}

	oldPassword, err := AesCBCDecrypt(oldPassword_iv, oldPassword, service.SecretService.Session.SymmetricKey)

	if err != nil {
		log.Panicf("Cannot decrypt old password. Error: %v", err)
	}

	newPassword, err = AesCBCDecrypt(newPassword_iv, newPassword, service.SecretService.Session.SymmetricKey)

	if err != nil {
		log.Panicf("Cannot decrypt new password. Error: %v", err)
	}

	oldSalt, err = AesCBCDecrypt(oldSalt_iv, oldSalt, service.SecretService.Session.SymmetricKey)

	if err != nil {
		log.Panicf("Cannot decrypt old salt. Error: %v", err)
	}

	newSalt, err = AesCBCDecrypt(newSalt_iv, newSalt, service.SecretService.Session.SymmetricKey)

	if err != nil {
		log.Panicf("Cannot decrypt new salt. Error: %v", err)
	}

	// set new password (requires previous empty password)
	if len(string(oldPassword[:])) < 1 {
		if len(service.ReadPasswordFile()) < 1 {
			hasher := sha512.New()
			hasher.Write(append(newSalt[:], newPassword[:]...))
			hash := hex.EncodeToString(hasher.Sum(nil))
			err = service.WritePasswordFile(hash)

			if err != nil {
				log.Panicf("Cannot write password file. Error: %v", err)
			}

		} else {
			return "password is not empty", nil
		}
	} else { // change old password
		// check if password match
		hasher := sha512.New()
		hasher.Write(append(oldSalt[:], oldPassword[:]...))
		hash := hex.EncodeToString(hasher.Sum(nil))
		if service.ReadPasswordFile() != hash {
			log.Warnf("Password mismatch from CLI client")
			return "wrong old password", nil
		}

		// old password matches, change password file
		hasher.Write(append(newSalt[:], newPassword[:]...))
		hash = hex.EncodeToString(hasher.Sum(nil))
		err := service.WritePasswordFile(hash)

		if err != nil {
			log.Panicf("Failed to write password file. Error: %v", err)
		}
	}

	return "ok", nil
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< SetPassword <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Command >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Command ( IN String command,
	          OUT String result);
*/

// Command receives a command from CLI and runs it on daemon side
func (service *Service) Command(
	serialnumber string,
	cookie []byte, cookie_iv []byte,
	command []byte, command_iv []byte,
	params []byte, params_iv []byte) ([]byte, []byte, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":    "ir.remisa.SecretService",
		"method":       "Command",
		"serialnumber": serialnumber,
		"cookie":       cookie,
		"cookie_iv":    cookie_iv,
		"command":      command,
		"command_iv":   command_iv,
		"params":       params,
		"params_iv":    params_iv,
	}).Trace("Method called by client")

	// TODO: Implement

	return nil, nil, nil

}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Command <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
