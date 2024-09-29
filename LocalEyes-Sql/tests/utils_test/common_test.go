package utils_test

import (
	"github.com/stretchr/testify/assert"
	"localEyes/utils"
	"testing"
)

func TestGenerateRandomId(t *testing.T) {
	// Calling the function multiple times to check for uniqueness
	ids := make(map[string]int)

	for i := 0; i < 10; i++ {
		id := utils.GenerateRandomId()
		assert.Len(t, id, 4, "Generated ID should be 4 characters long")
		ids[id] = 0 // Store unique IDs
	}

	assert.Equal(t, len(ids), 10, "Generated IDs should be unique")
}
