package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// MD5Hash returns the uppercase MD5 hash of the input string
func MD5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// SHA256Hash returns the uppercase SHA256 hash of the input string
func SHA256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// SHA256HashBytes returns the raw SHA256 hash bytes
func SHA256HashBytes(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// GenerateCnonce generates a random 8-character uppercase hex string
func GenerateCnonce() string {
	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to a fixed value in case of error (should not happen)
		return "ABCD1234"
	}
	return strings.ToUpper(hex.EncodeToString(bytes))
}
