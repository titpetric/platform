package model

import (
	"testing"
	"time"

	"github.com/titpetric/platform/internal/require"
)

func TestUser(t *testing.T) {
	m1 := NewUser()
	m1.FirstName = "Tit"
	m1.LastName = "Petric"

	m2 := NewUser()
	m2.FirstName = "Tit"
	m2.LastName = "Petric"
	m2.SetDeletedAt(time.Now())

	s1 := m1.String()
	s2 := m2.String()

	require.NotEqual(t, s1, s2)
	require.Equal(t, s1, "Tit Petric")
	require.Equal(t, s2, "Deleted user")
}
