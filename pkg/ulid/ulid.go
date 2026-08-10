package ulid

import (
	ulid "github.com/oklog/ulid/v2"
)

// ULID will return a new ulid.ULID value.
func ULID() ulid.ULID {
	return ulid.Make()
}

// String will return a string UUID.
func String() string {
	return ULID().String()
}

// Parse will parse a string into an ulid.ULID type.
// The returned value should be discarded in case of an error.
func Parse(in string) (ulid.ULID, error) {
	return ulid.Parse(in)
}

// Valid will return true if the string is a valid ULID.
func Valid(in string) bool {
	_, err := Parse(in)
	return err == nil
}
