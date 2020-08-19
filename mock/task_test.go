package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// runTest with random n times
const runTests = 10

// TestGetTaskMock checks random generation with ID 0
func TestGetMockTaskID0(t *testing.T) {
	for i := 0; i < runTests; i++ {
		tt := GetMockTaskID0()

		assert.Equal(t, uint(0), tt.ID)
		assert.NotEmpty(t, tt.Title)
		assert.NotEmpty(t, tt.Description)
		assert.NotNil(t, tt.Done)
	}
}

// TestGetTaskMock checks random generation
func TestGetTaskMock(t *testing.T) {
	for i := 0; i < runTests; i++ {
		tt := GetMockTask()

		assert.NotEqual(t, uint(0), tt.ID)
		assert.NotEmpty(t, tt.Title)
		assert.NotEmpty(t, tt.Description)
		assert.NotNil(t, tt.Done)
	}
}
