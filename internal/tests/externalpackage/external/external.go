package external

import (
	"github.com/gofrs/uuid"
)

// ProvidesUUID provides the uuid.UUID type
func ProvidesUUID() uuid.UUID {
	return uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000"))
}

// NeedsUUID needs the uuid type
func NeedsUUID(u uuid.UUID) bool {
	return true
}
