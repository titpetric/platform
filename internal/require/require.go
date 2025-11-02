package require

import "github.com/stretchr/testify/require"

var (
	New      = require.New
	True     = require.True
	False    = require.False
	Equal    = require.Equal
	NotEqual = require.NotEqual
	Error    = require.Error
	ErrorIs  = require.ErrorIs
	NoError  = require.NoError
	Nil      = require.Nil
	NotNil   = require.NotNil
)
