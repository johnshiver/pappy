package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const encryptionKeySize = 32

func hashAndSalt(pwd []byte) (string, error) {
	log.Println("hashing your pw...")
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("while hashing and salting password: %w", err)
	}
	return string(hash), nil
}

func comparePasswords(hashedPwd, plainPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	if err != nil {
		log.Debug(err)
		return false
	}
	return true
}

func generateEncryptionKey(components ...string) []byte {
	combinedComponents := strings.Join(components, "")
	// if not len encryptionKeySize, buffer until it is
	if len(combinedComponents) > encryptionKeySize {
		combinedComponents = combinedComponents[:encryptionKeySize]
	}
	for len(combinedComponents) < encryptionKeySize {
		combinedComponents += "d"
	}

	return []byte(combinedComponents)
}

// encrypt string to base64 crypto using AES
func encrypt(key, plaintext []byte) (string, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("while creating new cipher: %w", err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		 return "", fmt.Errorf("while reading ciphertext: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		log.Fatal("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
