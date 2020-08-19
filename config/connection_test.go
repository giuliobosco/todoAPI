package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTestInit test TestInit() db configurator
func TestTestInit(t *testing.T) {
	db := TestInit()

	assert.Equal(t, db, DB)
	assert.Equal(t, db, GetDB())
}
