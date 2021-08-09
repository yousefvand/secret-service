package client

import (
	"sync"

	"github.com/godbus/dbus/v5"
)

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Client >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// secret service client data structure
type Client struct {
	// dbus session connection
	Connection *dbus.Conn
	// dbus object used to call dbus methods
	DbusObject dbus.BusObject
	// Signal channel
	SignalChan chan *dbus.Signal
	// Mutex for lock/unlock Sessions map
	SessionsMutex *sync.RWMutex
	// sessions map. key: session dbus object path, value: session object
	Sessions map[string]*Session
	// Cli session
	CliSession *CliSession
	// Mutex for lock/unlock Collections map
	CollectionsMutex *sync.RWMutex
	// Collections map. key: Collection dbus object path, value: Collection object
	Collections map[string]*Collection
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Client <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Session >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// encryption algorithm type
// plain or Dh_ietf1024_sha256_aes128_cbc_pkcs7
type EncryptionAlgorithm uint8

const (
	// Plain algorithm (no encryption)
	Plain EncryptionAlgorithm = iota
	// Dh_ietf1024_sha256_aes128_cbc_pkcs7 algorithm
	Dh_ietf1024_sha256_aes128_cbc_pkcs7
	// Unsupported algorithm (used in tests)
	Unsupported
)

// Session data structure
type Session struct {
	// reference to parent (client)
	Parent *Client
	// session full dbus object path
	ObjectPath dbus.ObjectPath
	// encryption algorithm type
	EncryptionAlgorithm EncryptionAlgorithm
	// symmetric key used or AES encryption/decryption. Needs IV as well
	SymmetricKey []byte // 16 bytes (128 bits)
	// client public key used or AES encryption/decryption
	ServicePublicKey []byte // 128 bytes (1024 bits)
}

type CliSession struct {
	// reference to parent (client)
	Parent *Client
	// symmetric key used or AES encryption/decryption. Needs IV as well
	SymmetricKey []byte // 16 bytes (128 bits)
	// cookie
	Cookie string
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Session <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Collection >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

/*

Collection
├── item = secret + lookup attributes + label
├── item = secret + lookup attributes + label
├── item = secret + lookup attributes + label
└──
...

*/

// Collection data structure
// collection consists of items
type Collection struct {
	// reference to parent (client)
	Parent *Client
	// Signal channel
	SignalChan chan *dbus.Signal
	// Mutex for lock/unlock Items slice
	ItemsMutex *sync.RWMutex
	// Items map. key: Item dbus object path, value: Item object
	Items map[string]*Item
	// collection full dbus object path
	ObjectPath dbus.ObjectPath
	// collection rawProperties map
	rawProperties map[string]dbus.Variant
	// Mutex for lock/unlock Properties map
	PropertiesMutex *sync.RWMutex
	// collection Properties map
	Properties map[string]dbus.Variant
	// collection alias (friendly name)
	Alias string
	// Mutex to lock/unlock Locked status of collection
	LockMutex *sync.Mutex
	// collection Label
	Label string
	// true if collection is locked otherwise false
	Locked bool
	// Unix time collection created
	Created uint64
	// Unix time collection modified
	Modified uint64
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Collection <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Item >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// Item data structure
// item = secret + lookup attributes + label
type Item struct {
	// reference to parent (collection)
	Parent *Collection
	// item full dbus object path
	ObjectPath dbus.ObjectPath
	// secret contained in this item
	Secret *Secret
	// Mutex for lock/unlock LookupAttributes slice
	LookupAttributesMutex *sync.RWMutex
	// LookupAttributes (name + value) contained in this item
	LookupAttributes map[string]string
	// label of this item
	Label string
	// Mutex to lock/unlock Locked status of item
	LockMutex *sync.Mutex
	// true if item is locked otherwise false
	Locked bool
	// Unix time item created
	Created uint64
	// Unix time item modified
	Modified uint64
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Item <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Secret >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// Secret data structure
type Secret struct {
	// reference to parent (item)
	Parent *Item
	// Unencrypted secret
	PlainSecret string
	// Secret type needed bu API
	SecretApi *SecretApi
}

// Secret data structure needed bu API
type SecretApi struct {
	// The session full dbus object path that was used to encode the secret
	Session dbus.ObjectPath
	// Algorithm dependent parameters for secret value encoding
	Parameters []byte
	// Possibly encoded secret value
	Value []byte
	//The content type of the secret i.e. ‘text/plain; charset=utf8’
	ContentType string
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Secret <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Prompt >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// Prompt data structure
type Prompt struct {
	// reference to parent (client)
	Parent *Client
	// Signal channel
	SignalChan chan *dbus.Signal
	// prompt full dbus object path
	ObjectPath dbus.ObjectPath
	// client applications can use the window-id to
	// display the prompt attached to their application window
	WindowId string
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Prompt <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Secret Map >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// A mapping from object-paths to Secret structs
type SecretMap map[dbus.ObjectPath]Secret

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Secret Map <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
