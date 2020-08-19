package tu

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

//https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	bytes, err := GenerateRandomBytes(n)

	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes), nil
}

// GenerateRandomEmail generate an random email
func GenerateRandomEmail() (string, error) {
	name, err := GenerateRandomString(8)

	if err != nil {
		return "", err
	}

	domain, err := GenerateRandomString(6)

	if err != nil {
		return "", err
	}

	return name + "@" + domain + ".com", nil
}

// RandomEmail genrate a ranodm email ignoring errors
func RandomEmail() string {
	e, _ := GenerateRandomEmail()

	return e
}

// RandomString12 generate random string of length 12 ignoring errors
func RandomString12() string {
	s, _ := GenerateRandomString(12)
	return s
}

// RandomUint generate random int ingoring errors
func RandomUint() uint {
	b, _ := GenerateRandomBytes(1)

	return uint(b[0])
}

// RandomUintNo0 generate random int ignoring errors 0 not allowed
func RandomUintNo0() uint {
	for b, _ := GenerateRandomBytes(1); b[0] == 0; {
		return uint(b[0])
	}
	return 1
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

// PasswordHash creates a password hash of the string given as parameter
// For check the created hash is required to use ComparePasswordHash
// DO NOT compare manually 2 hashes.
func PasswordHash(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(passwordBytes), err
}

// ComparePasswordHash Compare password with Password Hash
func ComparePasswordHash(h string, p string) bool {
	hb := []byte(h)
	pb := []byte(p)

	err := bcrypt.CompareHashAndPassword(hb, pb)
	if err != nil {
		return false
	}

	return true
}
