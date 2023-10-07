package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"

	"github.com/ellofae/payment-system-kafka/config"
)

/*
Enctypion/Decryption Algorithm - Advanced Encryption Standard (AES) (AES Encryption)
*/

var encryptionKey string
var once sync.Once

func InitializeEncryptionKey(cfg *config.Config) {
	once.Do(func() {
		encryptionKey = cfg.Encryption.EncryptionKey
	})
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EncryptData(data []byte) string {
	block, _ := aes.NewCipher([]byte(createHash(encryptionKey)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return string(ciphertext)
}

func DecryptData(data []byte) string {
	key := []byte(createHash(encryptionKey))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintext)
}
