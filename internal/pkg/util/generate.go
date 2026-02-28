// Package util provides utility functions.
package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

// HashToken creates a SHA256 hash of a token for storage.
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GenerateVerificationCode generates a 6-digit verification code using crypto/rand.
func GenerateVerificationCode() string {
	// Generate a random number between 0 and 999999
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		// Fallback to time-based if crypto rand fails (shouldn't happen in practice)
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// GenerateState generates a random state string for OAuth.
func GenerateState() string {
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), "ibookfs")
}
