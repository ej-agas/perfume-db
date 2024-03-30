package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

func PKCS7Pad(data []byte, blockSize int) []byte {
	padSize := blockSize - (len(data) % blockSize)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(data, pad...)
}

func PKCS7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("pkcs7: invalid padding")
	}
	padSize := int(data[len(data)-1])
	if padSize > len(data) {
		return nil, errors.New("pkcs7: invalid padding")
	}
	return data[:len(data)-padSize], nil
}

func (app *application) Encrypt(plaintext []byte) (string, error) {
	// Create a new AES block cipher with the provided key
	block, err := aes.NewCipher([]byte(app.config.encryptionKey))
	if err != nil {
		return "", err
	}

	plaintext = PKCS7Pad(plaintext, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	encoded := base64.URLEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

func (app *application) Decrypt(encoded string) ([]byte, error) {
	// Create a new AES block cipher with the provided key
	ciphertext, err := base64.URLEncoding.DecodeString(encoded)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(app.config.encryptionKey))
	if err != nil {
		return nil, err
	}

	// Extract the IV from the beginning of the ciphertext

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("invalid cipher text")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Create a new AES cipher block chaining (CBC) mode
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext, err = PKCS7Unpad(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
