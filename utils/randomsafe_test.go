package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	repeetForTimes = 20 // Repeet Test for 20 random times.
	stringLength   = 12 // Length of the generated strings
)

func TestGenerateRandomBytes(t *testing.T) {
	for i := 0; i < repeetForTimes; i++ {
		b, err := GenerateRandomBytes(i)

		assert.Equal(t, nil, err)
		assert.Equal(t, i, len(b))
	}
}

func TestPasswordHashAndComparation(t *testing.T) {
	for i := 0; i < repeetForTimes; i++ {
		p, err := GenerateRandomString(stringLength)

		assert.Equal(t, nil, err)

		h, err := PasswordHash(p)

		assert.Equal(t, nil, err)

		for j := 0; j < repeetForTimes; j++ {
			c := ComparePasswordHash(h, p)

			assert.True(t, c)
		}
	}
}
