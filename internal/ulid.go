package internal

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func ULID() string {
	now := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(now), entropy).String()
}
