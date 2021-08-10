// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
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
	OpenSession ( IN String algorithm,
	              IN Variant input,
	              OUT Variant output,
	              OUT String serialnumber);
*/

// OpenSession opens a unique session for the caller application
// further communication encryption/decryption relies on the related session
func (secretservice *SecretService) OpenSession(algorithm string,
	input dbus.Variant) (dbus.Variant, string, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface": "ir.remisa.SecretService",
		"method":    "OpenSession",
		"algorithm": algorithm,
		"input":     input.Value(),
	}).Trace("Method called by client")

	log.Debugf("Client suggested '%s' algorithm", algorithm)

	// if OpenSession succeeds all related information are stored in a session
	cliSession := NewCliSession(secretservice.Parent)
	secretservice.Parent.CliSession = cliSession

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
		cliSession.SymmetricKey = symmetricKey

		log.Tracef("Symmetric key: %v", cliSession.SymmetricKey)
		log.Tracef("Symmetric key length: %v", len(cliSession.SymmetricKey))

		log.Debug("Agreed on 'dh-ietf1024-sha256-aes128-cbc-pkcs7' algorithm")

		serialnumber := UUID()

		return dbus.MakeVariant(publicKey), serialnumber, nil // end of successful negotiation

	default: // algorithm is not 'plain' or 'dh-ietf1024-sha256-aes128-cbc-pkcs7'
		log.Warnf("The '%s' algorithm suggested by client is not supported", algorithm)
		return dbus.MakeVariant(""), "", ApiErrorNotSupported()

	}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< OpenSession <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Command >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*
	Command ( IN String command,
	          OUT String result);
*/

// Command receives a command from CLI and runs it on daemon side
func (secretservice *SecretService) Command(
	serialnumber string, command dbus.Variant, params []dbus.Variant) ([]dbus.Variant, *dbus.Error) {

	log.WithFields(log.Fields{
		"interface":    "ir.remisa.SecretService",
		"method":       "Command",
		"serialnumber": serialnumber,
		"command":      command,
		"params":       params,
	}).Trace("Method called by client")

	return nil, nil

}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Command <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
