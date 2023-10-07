package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"time"
)

func GenerateUniqueRandomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	timestamp := time.Now().UnixNano()
	randomString += "_" + strconv.FormatInt(timestamp, 10)

	return randomString
}
