package uuid

import "github.com/google/uuid"

type UUID interface {
	string | []byte | uuid.UUID
}
