# GoKit UUID (guid)

[![Go Reference](https://pkg.go.dev/badge/github.com/gokit/guid.svg)](https://pkg.go.dev/github.com/gokit/guid)

A lightweight wrapper around Google's UUID library that provides convenient functions for handling and converting between various UUID representations.

## Overview

The `guid` package offers utility functions for working with UUIDs (Universally Unique Identifiers) in different forms. It simplifies common operations like converting between string, binary, and UUID object representations while handling error cases gracefully.

## Installation

```bash
go get github.com/gokit/guid
```

## Usage

```go
import (
    "fmt"
    "github.com/gokit/guid"
    "github.com/google/uuid"
)

func main() {
    // Create a new UUID
    originalUUID := uuid.New()
    fmt.Printf("Original UUID: %s\n", originalUUID.String())

    // Convert UUID to string and back
    uuidStr := originalUUID.String()
    parsedUUID := guid.UUIDFromString(uuidStr)
    fmt.Printf("Parsed UUID from string: %s\n", parsedUUID.String())

    // Convert UUID to binary representation
    binaryUUID := guid.ByteSliceFromUUID(originalUUID)
    fmt.Printf("Binary UUID length: %d bytes\n", len(binaryUUID))

    // Convert binary back to UUID object
    fromBinaryUUID := guid.UUIDFromByteSlice(binaryUUID)
    fmt.Printf("UUID from binary: %s\n", fromBinaryUUID.String())

    // Convert string UUID directly to binary
    binaryFromStr := guid.ByteSliceFromStringUUID(uuidStr)

    // Convert binary UUID directly to string
    strFromBinary := guid.StringFromByteSliceUUID(binaryUUID)
    fmt.Printf("String from binary: %s\n", strFromBinary)
}
```

## Functions

The package provides the following functions:

- `UUIDFromString(s string) uuid.UUID`: Converts a string to a UUID object
- `UUIDFromByteSlice(b []byte) uuid.UUID`: Creates a UUID from a byte slice
- `ByteSliceFromStringUUID(s string) []byte`: Converts a UUID string to its binary representation
- `StringFromByteSliceUUID(b []byte) string`: Converts a binary UUID to its string form
- `ByteSliceFromUUID(u uuid.UUID) []byte`: Converts a UUID object to its binary representation

All functions have proper error handling - they return zero values (nil, empty string, or uuid.Nil) when given invalid input.

## Error Handling

All functions in this package handle errors internally and return appropriate zero values when encountering issues:

- Invalid UUID strings return `uuid.Nil` or `nil` byte slices
- Invalid byte slices return `uuid.Nil` or its string representation
- Marshal errors return `nil` byte slices

## Dependencies

- [github.com/google/uuid](https://github.com/google/uuid): The underlying UUID implementation

## License

MIT License

Copyright (c) 2023, The GoKit Authors
