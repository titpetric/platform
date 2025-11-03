package assert

import "github.com/stretchr/testify/assert"

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
)
