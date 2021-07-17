package client

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"strconv"

	"github.com/godbus/dbus/v5"
	"github.com/monnand/dhkx"
	"golang.org/x/crypto/hkdf"
)

/*
	OpenSession ( IN String algorithm,
	              IN Variant input,
	              OUT Variant output,
	              OUT ObjectPath result);
*/

// OpenSession creates a session for encrypted or non-encrypted further communication
func (client *Client) OpenSession(algorithm EncryptionAlgorithm) (*Session, error) {

	var algorithmInUse string
	var input dbus.Variant
	session := NewSession(client)
	session.EncryptionAlgorithm = algorithm

	group, errGroup := dhkx.GetGroup(2) // 2 -> 1024 bit (128 bytes) secret key
	if errGroup != nil {
		return nil, errors.New("Diffie–Hellman group creation failed. Error: " + errGroup.Error())
	}

	privateKey, errPrivateKey := group.GeneratePrivateKey(rand.Reader)
	if errPrivateKey != nil {
		return nil,
			errors.New("Diffie–Hellman private key generation failed. Error: " + errPrivateKey.Error())
	}

	switch algorithm {

	case Plain:
		algorithmInUse = "plain"
		input = dbus.MakeVariant("")

	case Dh_ietf1024_sha256_aes128_cbc_pkcs7:
		algorithmInUse = "dh-ietf1024-sha256-aes128-cbc-pkcs7"
		input = dbus.MakeVariant(privateKey.Bytes()) // own public key

	default: // Unsupported (used in tests)
		algorithmInUse = "unsupported"
		input = dbus.MakeVariant("")
	}

	var call *dbus.Call
	var err error

	call, err = client.Call("org.freedesktop.secrets", "/org/freedesktop/secrets",
		"org.freedesktop.Secret.Service", "OpenSession", algorithmInUse, input)

	if err != nil {
		return nil, errors.New("dbus call failed. Error: " + err.Error())
	}

	var output dbus.Variant
	var result dbus.ObjectPath
	err = call.Store(&output, &result)

	if err != nil {
		if algorithm == Unsupported {
			return nil, errors.New("unsupported encryption algorithm")
		} else {
			return nil, errors.New("type conversion failed in 'OpenSession'. Error: " + err.Error())
		}
	}

	if algorithm == Dh_ietf1024_sha256_aes128_cbc_pkcs7 {
		var servicePublicKey []byte
		err = dbus.Store([]interface{}{output.Value()}, &servicePublicKey)
		if err != nil {
			return nil, errors.New("Cannot convert client public key. Error: " + err.Error())
		}
		if len(servicePublicKey) != 128 {
			return nil, errors.New("Wrong length of public key. Expected 128 bytes got " +
				strconv.Itoa(len(servicePublicKey)) + " bytes")
		}
		session.ServicePublicKey = servicePublicKey // dhkx.NewPublicKey(clientPublicKey)

		sharedKey, err := group.ComputeKey(dhkx.NewPublicKey(servicePublicKey), privateKey)
		if err != nil {
			return nil,
				errors.New("Diffie–Hellman shared key generation failed. Error: " + err.Error())
		}
		sessionSharedKey := sharedKey.Bytes()

		hkdf := hkdf.New(sha256.New, sessionSharedKey, nil, nil)
		symmetricKey := make([]byte, aes.BlockSize) // 16 * 8 = 128 bit
		n, err := io.ReadFull(hkdf, symmetricKey)
		if n != aes.BlockSize {
			return nil,
				errors.New("Cannot create 16 byte key. Length is: " + strconv.Itoa(len(symmetricKey)))
		}
		if err != nil {
			return nil, errors.New("Symmetric Key generation failed. Error: " + err.Error())
		}
		session.SymmetricKey = symmetricKey

	}

	session.ObjectPath = result
	client.AddSession(session)

	return session, nil
}
