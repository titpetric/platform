// Package require provides test requirement helpers re-exported from testify/require.
package require

import "github.com/stretchr/testify/require"

// Exported requirement functions from testify/require.
var (
	New           = require.New
	True          = require.True
	False         = require.False
	Equal         = require.Equal
	NotEqual      = require.NotEqual
	Error         = require.Error
	ErrorIs       = require.ErrorIs
	ErrorContains = require.ErrorContains
	NoError       = require.NoError
	Nil           = require.Nil
	NotNil        = require.NotNil
	Empty         = require.Empty
	NotEmpty      = require.NotEmpty
)
