package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// UUID returns uuid without dashes
func UUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// Path2Name takes a dbus path like /a/b/c/d and returns
// a dbus object name like a.b.c and last part d
// Example: path=/a/b/c/xyz, name=Foo -> (a.b.c.Foo, xyz)
func Path2Name(path string, name string) (string, string) {

	path = strings.TrimSpace(path)
	name = strings.TrimSpace(name)

	result := strings.ReplaceAll(path, "/", ".")
	result = strings.TrimLeft(result, ".")
	idx := strings.LastIndex(result, ".")
	child := result[idx+1:]
	result = result[:idx]
	result += "." + name

	return result, child
}

// MemUsageOS returns memory usage on OS in MB
func MemUsageOS() uint64 {
	var memoryUsage runtime.MemStats
	runtime.ReadMemStats(&memoryUsage)
	return memoryUsage.Sys / 1024 / 1024
}

// IsMapSubsetSingleMatch returns true if only one
// key/value of mapSubset exists in mapSet otherwise false
func IsMapSubsetSingleMatch(mapSet map[string]string,
	mapSubset map[string]string, lock *sync.RWMutex) bool {

	lock.RLock()
	defer lock.RUnlock()

	if len(mapSubset) == 0 {
		return true
	}

	if len(mapSubset) > len(mapSet) {
		return false
	}

	for k, v := range mapSubset {
		if _, ok := mapSet[k]; ok {
			if mapSet[k] == v {
				return true
			}
		}
	}
	return false
}

// IsMapSubsetFullMatch returns true if mapSubset is a full subset of mapSet otherwise false
func IsMapSubsetFullMatch(mapSet map[string]string,
	mapSubset map[string]string, lock *sync.RWMutex) bool {

	lock.RLock()
	defer lock.RUnlock()

	if len(mapSubset) == 0 {
		return true
	}

	if len(mapSubset) > len(mapSet) {
		return false
	}

	for k, v := range mapSubset {
		if _, ok := mapSet[k]; !ok {
			return false
		}
		if mapSet[k] != v {
			return false
		}
	}
	return true
}

// IsMapSubsetFullMatchGeneric returns true if
// mapSubset is a full subset of mapSet otherwise false
func IsMapSubsetFullMatchGeneric(mapSet interface{},
	mapSubset interface{}, lock *sync.RWMutex) bool {

	lock.RLock()
	defer lock.RUnlock()

	mapSetValue := reflect.ValueOf(mapSet)
	mapSubsetValue := reflect.ValueOf(mapSubset)

	if fmt.Sprintf("%T", mapSet) != fmt.Sprintf("%T", mapSubset) {
		return false
	}

	if len(mapSetValue.MapKeys()) < len(mapSubsetValue.MapKeys()) {
		return false
	}

	if len(mapSubsetValue.MapKeys()) == 0 {
		return true
	}

	iterMapSubset := mapSubsetValue.MapRange()

	for iterMapSubset.Next() {
		k := iterMapSubset.Key()
		v := iterMapSubset.Value()

		value := mapSetValue.MapIndex(k)

		if !value.IsValid() || v.Interface() != value.Interface() {
			return false
		}
	}

	return true
}

/* TODO: Needs go 1.18

func IsMapSubsetFullMatchGeneric[K, V comparable](m, sub map[K]V) bool {
    if len(sub) > len(m) {
        return false
    }
    for k, vsub := range sub {
        if vm, found := m[k]; !found || vm != vsub {
            return false
        }
    }
    return true
}

*/

// CommandExists returns true if command exists on OS otherwise false
func CommandExists(cmdName string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+cmdName)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Epoch returns seconds from midnight 1/1/1970
func Epoch() uint64 {
	return uint64(time.Now().Unix())
}

////////////////////////////// Crypto //////////////////////////////

// PKCS7Padding pads given data to match encryption block size
func PKCS7Padding(plainUnpaddedData []byte, blockSize int) []byte {
	paddingSize := blockSize - len(plainUnpaddedData)%blockSize
	padData := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(plainUnpaddedData, padData...)
}

// PKCS7Padding unpads given data to original form
func PKCS7UnPadding(plainPaddedData []byte) []byte {
	length := len(plainPaddedData)
	unpadding := int(plainPaddedData[length-1])
	return plainPaddedData[:(length - unpadding)]
}

// AesCBCEncrypt encrypts data and returns iv, cipherData, error
func AesCBCEncrypt(plainData, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicf("Cannot create cipher with key: '%v'. Error: %v", key, err)
	}

	blockSize := block.BlockSize() // 16
	plainData = PKCS7Padding(plainData, blockSize)

	cipherData := make([]byte, len(plainData))
	// Initial vector IV must be unique, but does not need to be kept secret
	iv := make([]byte, blockSize)
	// Fill iv with random bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Panicf("Cannot fill IV with random bytes. Error: %v", err)
	}

	// block size and initial vector size must be the same
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherData, plainData)

	return iv, cipherData, nil
}

// AesCBCDecrypt decrypts cipher to original. returns data, error
func AesCBCDecrypt(iv, cipherData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicf("Cannot create cipher with key: '%v'. Error: %v", key, err)
	}

	blockSize := block.BlockSize() // 16

	if len(cipherData) < blockSize {
		log.Errorf("cipher data (%d bytes) < blocksize (%d bytes)", len(cipherData), blockSize)
		return nil, errors.New("cipher data too short")
	}

	// CBC mode always works in whole blocks.
	if len(cipherData)%blockSize != 0 {
		log.Errorf("cipher data (%d bytes) is not a multiple of the block size (%d bytes)",
			len(cipherData), blockSize)
		return nil, errors.New("cipher data is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	var plainData = make([]byte, len(cipherData))
	mode.CryptBlocks(plainData, cipherData)
	plainData = PKCS7UnPadding(plainData)

	return plainData, nil
}

////////////////////////////// db encryption //////////////////////////////

func EncryptAESCBC256(key string, text string) (string, error) {
	rawKey := []byte(key)
	rawText := []byte(text)

	aesCipher, err := aes.NewCipher(rawKey)

	if err != nil {
		return "", fmt.Errorf("cannot create new cipher. Error: %v", err)
	}

	gcm, err := cipher.NewGCM(aesCipher)

	if err != nil {
		return "", fmt.Errorf("cannot create new GCM. Error: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("cannot read random bytes. Error: %v", err)
	}

	cipher := gcm.Seal(nonce, nonce, rawText, nil)
	cipherBase64 := base64.URLEncoding.EncodeToString(cipher)

	return cipherBase64, nil
}

func DecryptAESCBC256(key string, cipherText string) (string, error) {

	rawKey := []byte(key)
	textCipher, err := base64.URLEncoding.DecodeString(cipherText)

	if err != nil {
		return "", fmt.Errorf("cannot base64 decode cipher text. Error: %v", err)
	}

	rawCipher := []byte(textCipher)

	aesCipher, err := aes.NewCipher(rawKey)

	if err != nil {
		return "", fmt.Errorf("cannot create new cipher. Error: %v", err)
	}

	gcm, err := cipher.NewGCM(aesCipher)

	if err != nil {
		return "", fmt.Errorf("cannot create new GCM. Error: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(rawCipher) < nonceSize {
		return "", fmt.Errorf("cipher text smaller than nonce size")
	}

	nonce, ciphertext := rawCipher[:nonceSize], rawCipher[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return "", fmt.Errorf("GCM open failed: Err: %v", err)
	}

	return string(plaintext), nil

}
