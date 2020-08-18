package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetTesting(t *testing.T) {
	value := true
	isTesting = value

	assert.Equal(t, value, IsTesting())

	value = !value
	SetTesting(value)

	assert.Equal(t, value, IsTesting())

	SetTesting(!value)

	assert.Equal(t, value, IsTesting())
}
