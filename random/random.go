// Package random provides random number/string generation and hashing helpers.
package random

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"math/big"
	mathrand "math/rand"

	"github.com/Maximumsoft-Co-LTD/go-util/format"
	"github.com/Maximumsoft-Co-LTD/go-util/slice"
)

// UniqueInts generates numberOfRandom unique random integers in [0, numberOfRange).
// Returns nil if numberOfRandom > numberOfRange.
func UniqueInts(numberOfRandom int, numberOfRange int) []int {
	if numberOfRandom > numberOfRange || numberOfRandom <= 0 || numberOfRange <= 0 {
		return nil
	}

	result := make([]int, 0, numberOfRandom)
	for len(result) < numberOfRandom {
		n := mathrand.Intn(numberOfRange)
		if !slice.Contains(result, n) {
			result = append(result, n)
		}
	}
	return result
}

// String generates a random string of the specified length using alphanumeric characters.
func String(length int) (string, error) {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("generating random byte: %w", err)
		}
		result[i] = charset[n.Int64()]
	}
	return string(result), nil
}

// HashMD5 returns the MD5 hex digest of the input string.
// Note: MD5 is not cryptographically secure — use for checksums only, not security.
func HashMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return fmt.Sprintf("%x", hash)
}

// HashSHA1 returns the SHA1 hex digest of the input string.
// Note: SHA1 is not cryptographically secure — use for checksums only, not security.
func HashSHA1(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Decimal generates a random decimal in [0.01, 1.0) with the given precision digits.
// maxAttempts limits the retry loop to prevent hangs.
func Decimal(precision int, maxAttempts int) (float64, error) {
	if maxAttempts <= 0 {
		maxAttempts = 1000
	}
	for i := 0; i < maxAttempts; i++ {
		decimal := format.ToFixed(mathrand.Float64(), precision)
		if decimal >= 0.01 && decimal < 1 {
			return decimal, nil
		}
	}
	return 0, fmt.Errorf("failed to generate valid decimal after %d attempts", maxAttempts)
}
