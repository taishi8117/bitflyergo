package bitflyergo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func sign(message string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
