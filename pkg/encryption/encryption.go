package encryption

import (
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

func EncryptData(key []byte, text string) string {
	return ""
}

func DecryptData(key []byte, ciphertext string) string {
	return ""
}
