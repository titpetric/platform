package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	// "github.com/titpetric/platform/module/user/model"
)

// type User = model.User
func TestUser_String(t *testing.T) {
	m1 := &User{
		FirstName: "Tit",
		LastName:  "Petric",
	}

	m2 := &User{
		FirstName: "Tit",
		LastName:  "Petric",
	}
	m2.SetDeletedAt(time.Now())

	s1 := m1.String()
	s2 := m2.String()

	require.NotEqual(t, s1, s2)
	require.Equal(t, s1, "Tit Petric")
	require.Equal(t, s2, "Deleted user")
}
