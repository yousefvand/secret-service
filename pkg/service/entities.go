// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

import (
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Service >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// SaveData is a function used by a child
// to inform parent of a change in data
type SaveData func()

// secretservice data structure
type Service struct {
	// dbus session connection
	Connection *dbus.Conn

	// TODO: Remove
	// Mutex for lock/unlock service
	// LockMutex *sync.Mutex
	// true if service is locked otherwise false
	// Locked bool

	// Service home path
	Home string
	// encrypt database
	EncryptDatabase bool
	// SecretService session
	SecretService *SecretService
	// Mutex for lock/unlock Sessions map
	SessionsMutex *sync.RWMutex
	// Cli Session
	CliSession *CliSession // TODO: REMOVE ME
	// sessions map. key: session dbus object path, value: session object
	Sessions map[string]*Session
	// Mutex for lock/unlock Collections map
	CollectionsMutex *sync.RWMutex
	// Collections map. key: Collection dbus object path, value: Collection object
	Collections map[string]*Collection
	// inform parent data has happened
	// SaveData SaveData
	// Channel to signal saving data to db
	SaveSignalChan chan struct{}
	// inform service is up and ready
	ServiceReadyChan chan struct{}
	// inform service is shutdown
	ServiceShutdownChan chan struct{}
	// inform database has loaded
	DbLoadedChan chan struct{}
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Service <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> SecretService >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// CLI interface data structure
type SecretService struct {
	// reference to parent (service)
	Parent *Service
	// session (public key negotiation)
	Session *SecretServiceCLiSession
}

// session (public key negotiation)
type SecretServiceCLiSession struct {
	//  session serial number
	SerialNumber string
	// symmetric key used or AES encryption/decryption. Needs IV as well
	SymmetricKey []byte // 16 bytes (128 bits)
	// session cookie
	Cookie Cookie
}

type Cookie struct {
	Value  string
	Issued time.Time
	time.Duration
}

type PasswordFile struct {
	// Password file version
	Version string `yaml:"version"`
	// Password hash: sha512(salt+password)
	PasswordHash string `yaml:"passwordHash"`
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< SecretService <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Session >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// encryption algorithm type
// plain or Dh_ietf1024_sha256_aes128_cbc_pkcs7
type EncryptionAlgorithm uint8

const (
	// Plain algorithm (no encryption)
	Plain EncryptionAlgorithm = iota
	// Dh_ietf1024_sha256_aes128_cbc_pkcs7 algorithm
	Dh_ietf1024_sha256_aes128_cbc_pkcs7
)

// Session data structure
type Session struct {
	// reference to parent (service)
	Parent *Service
	// session full dbus object path
	ObjectPath dbus.ObjectPath
	// encryption algorithm type
	EncryptionAlgorithm EncryptionAlgorithm
	// symmetric key used or AES encryption/decryption. Needs IV as well
	SymmetricKey []byte // 16 bytes (128 bits)
	// Sessions don't need to get persistent in db so no need for 'Update'
}

type CliSession struct {
	// reference to parent (service)
	Parent *Service
	// symmetric key used or AES encryption/decryption. Needs IV as well
	SymmetricKey []byte // 16 bytes (128 bits)
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
	// reference to parent (service)
	Parent *Service
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
	// dbus properties handle
	DbusProperties *prop.Properties
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
	// inform parent data has happened
	SaveData SaveData

	// Temporary solution to data race in marshaling for db
	DataMutex *sync.RWMutex
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
	// Mutex for lock/unlock Properties map
	PropertiesMutex *sync.RWMutex
	// collection Properties map
	Properties map[string]dbus.Variant
	// item properties map
	rawProperties map[string]dbus.Variant
	// dbus properties handle
	DbusProperties *prop.Properties
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
	// inform parent data has happened
	SaveData SaveData

	// Temporary solution to data race in marshaling for db
	DataMutex *sync.RWMutex
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Item <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Secret >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// Secret data structure
type Secret struct {
	// reference to parent (item)
	Parent *Item
	// Unencrypted secret
	PlainSecret string
	// Secret type needed by API
	SecretApi *SecretApi
	// inform parent data has happened
	SaveData SaveData

	// Temporary solution to data race in marshaling for db
	DataMutex *sync.RWMutex
}

// Secret data structure needed bu API
type SecretApi struct {
	// The session full dbus object path that was used to encode the secret
	Session dbus.ObjectPath
	// Algorithm dependent parameters for secret value encoding
	Parameters []byte
	// Possibly encoded secret value
	Value []byte
	//The content type of the secret i.e. 'text/plain; charset=utf8'
	ContentType string
}

/* <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Secret <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */

/* >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Prompt >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */

// Prompt data structure
type Prompt struct {
	// reference to parent
	Parent *Service
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
