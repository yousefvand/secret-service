package client

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"
)

// SliceContains returns true if a slice contains an element otherwise false
func SliceContains(slice, elem interface{}) (bool, error) {

	sv := reflect.ValueOf(slice)

	// Check that slice is actually a slice/array.
	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Array {
		return false, errors.New("not an array or slice")
	}

	// iterate the slice
	for i := 0; i < sv.Len(); i++ {

		// compare elem to the current slice element
		if elem == sv.Index(i).Interface() {
			return true, nil
		}
	}

	// nothing found
	return false, nil

}

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
		panic(fmt.Sprintf("Cannot create cipher with key: '%v'. Error: %v", key, err))
	}

	blockSize := block.BlockSize() // 16
	plainData = PKCS7Padding(plainData, blockSize)

	cipherData := make([]byte, len(plainData))
	// Initial vector IV must be unique, but does not need to be kept secret
	iv := make([]byte, blockSize)
	// Fill iv with random bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(fmt.Sprintf("Cannot fill IV with random bytes. Error: %v", err))
	}

	// block size and initial vector size must be the same
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherData, plainData)

	return iv, cipherData, nil
}

// AesCBCDecrypt decrypts cipher to original data
func AesCBCDecrypt(iv, cipherData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("Cannot create cipher with key: '%v'. Error: %v", key, err))
	}

	blockSize := block.BlockSize() // 16

	if len(cipherData) < blockSize {
		return nil, errors.New("cipher data too short")
	}

	// CBC mode always works in whole blocks.
	if len(cipherData)%blockSize != 0 {
		return nil, errors.New("cipher data is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	var plainData = make([]byte, len(cipherData))
	mode.CryptBlocks(plainData, cipherData)
	plainData = PKCS7UnPadding(plainData)

	return plainData, nil
}
