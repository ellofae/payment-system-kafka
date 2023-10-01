package encryption

import (
	"crypto/aes"
	"encoding/hex"
	"sync"

	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/pkg/logger"
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

func EncryptData(data string) (string, error) {
	log := logger.GetLogger()

	c, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		log.Error("Unable to create a Cipher block", "error", err.Error())
		return "", err
	}

	out := make([]byte, len(data))

	c.Encrypt(out, []byte(data))

	return hex.EncodeToString(out), nil
}

func DecryptData(ct string) (string, error) {
	log := logger.GetLogger()

	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		log.Error("Unable to create a Cipher block", "error", err.Error())
		return "", err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s, nil
}
