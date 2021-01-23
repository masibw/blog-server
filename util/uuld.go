package util

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func Generate(time time.Time) string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.UnixNano())), 0) // nolint:gosec
	return ulid.MustNew(ulid.Timestamp(time), entropy).String()
}
