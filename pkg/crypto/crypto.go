package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	// log "github.com/sirupsen/logrus"
)

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
		panic("Cannot create cipher with key: '" + string(key) + "'. Error: " + err.Error())
	}

	blockSize := block.BlockSize() // 16
	plainData = PKCS7Padding(plainData, blockSize)

	cipherData := make([]byte, len(plainData))
	// Initial vector IV must be unique, but does not need to be kept secret
	iv := make([]byte, blockSize)
	// Fill iv with random bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic("Cannot fill IV with random bytes. Error: " + err.Error())
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
		panic("Cannot create cipher with key: '" + string(key) + "'. Error: " + err.Error())
	}

	blockSize := block.BlockSize() // 16

	if len(cipherData) < blockSize {
		// "cipher data (%d bytes) < blocksize (%d bytes)", len(cipherData), blockSize
		return nil, errors.New("cipher data too short")
	}

	// CBC mode always works in whole blocks.
	if len(cipherData)%blockSize != 0 {
		// "cipher data (%d bytes) is not a multiple of the block size (%d bytes)", len(cipherData), blockSize
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
