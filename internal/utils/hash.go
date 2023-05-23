package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MakeHashFromAdress returns md5 representation of given string by len of 5.
// It makes 62^5 potential combinations of string and a bit saves from collision in a low amount of requests but makes url short.
func MakeHashFromAdress(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])[:5]
}
