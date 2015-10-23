package lib

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// HMACBytes returns HMAC as bytes
func HMACBytes(h func() hash.Hash, s, k []byte) string {
	mac := hmac.New(h, k)
	mac.Write(s)
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// HMAC returns HMAC as string
func HMAC(h func() hash.Hash, s, k string) string {
	return string(HMACBytes(h, []byte(s), []byte(k)))
}

// HMACSha1Bytes returns sha1 encrypted HMAC as bytes
func HMACSha1Bytes(s, k []byte) string {
	return HMACBytes(sha1.New, s, k)
}

// HMACSha1 returns sha1 encrypted HMAC as string
func HMACSha1(s, k string) string {
	return HMAC(sha1.New, s, k)
}

// HMACSha256Bytes returns sha256 encrypted HMAC as bytes
func HMACSha256Bytes(s, k []byte) string {
	return HMACBytes(sha256.New, s, k)
}

// HMACSha256 returns sha256 encrypted HMAC as bytes
func HMACSha256(s, k string) string {
	return HMAC(sha256.New, s, k)
}

// Sha1Bytes calculates Sha1 hash
func Sha1Bytes(s []byte) string {
	h := sha1.New()
	h.Write(s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
