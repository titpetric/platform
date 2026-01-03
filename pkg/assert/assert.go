// Package assert provides test assertion helpers re-exported from testify/assert.
package assert

import "github.com/stretchr/testify/assert"

// Exported assertion functions from testify/assert.
var (
	New           = assert.New
	True          = assert.True
	False         = assert.False
	Equal         = assert.Equal
	NotEqual      = assert.NotEqual
	Error         = assert.Error
	ErrorIs       = assert.ErrorIs
	ErrorContains = assert.ErrorContains
	NoError       = assert.NoError
	Nil           = assert.Nil
	NotNil        = assert.NotNil
	Empty         = assert.Empty
	NotEmpty      = assert.NotEmpty
)
