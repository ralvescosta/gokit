// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package guid provides utility functions for working with UUIDs.
// This package wraps the github.com/google/uuid library to provide
// convenient conversion functions between various UUID representations.
package guid

import "github.com/google/uuid"

// UUIDFromString parses a UUID string and returns a UUID object.
// If the string is not a valid UUID, it returns uuid.Nil.
func UUIDFromString(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}

	return u
}

// UUIDFromByteSlice creates a UUID from a byte slice.
// If the byte slice is not a valid UUID representation, it returns uuid.Nil.
func UUIDFromByteSlice(b []byte) uuid.UUID {
	u, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil
	}

	return u
}

// ByteSliceFromStringUUID converts a UUID string to its binary representation.
// If the string is not a valid UUID, it returns nil.
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

// StringFromByteSliceUUID converts a binary UUID representation to its string form.
// If the byte slice is not a valid UUID, it returns the string representation of uuid.Nil.
func StringFromByteSliceUUID(b []byte) string {
	u, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil.String()
	}

	return u.String()
}

// ByteSliceFromUUID converts a UUID object to its binary representation.
// If marshaling fails, it returns nil.
func ByteSliceFromUUID(u uuid.UUID) []byte {
	b, err := u.MarshalBinary()
	if err != nil {
		return nil
	}

	return b
}
