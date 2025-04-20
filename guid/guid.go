// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package guid

import "github.com/google/uuid"

func UUIDFromString(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}

	return u
}

func UUIDFromByteSlice(b []byte) uuid.UUID {
	u, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil
	}

	return u
}

func ByteSliceFromStringUUID(s string) []byte {
	u, err := uuid.Parse(s)
	if err != nil {
		return nil
	}

	b, err := u.MarshalBinary()
	if err != nil {
		return nil
	}

	return b
}

func StringFromByteSliceUUID(b []byte) string {
	u, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil.String()
	}

	return u.String()
}

func ByteSliceFromUUID(u uuid.UUID) []byte {
	b, err := u.MarshalBinary()
	if err != nil {
		return nil
	}

	return b
}
