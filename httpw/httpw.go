// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package httpw

// Message formats a message with the httpw package prefix.
func Message(msg string) string {
	return "[gokit:httpw] " + msg
}
