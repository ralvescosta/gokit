// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package mqtt

// QoS represents the Quality of Service levels for MQTT messages.
type QoS byte

const (
	// AtMostOnce represents QoS level 0.
	AtMostOnce QoS = 0
	// AtLeastOnce represents QoS level 1.
	AtLeastOnce QoS = 1
	// ExactlyOnce represents QoS level 2.
	ExactlyOnce QoS = 2
)

// LogMessage formats and returns a log message with a consistent prefix.
func LogMessage(msg ...string) string {
	prefix := "[gokit::mqtt] "
	for _, s := range msg {
		prefix += s
	}
	return prefix
}

// QoSFromBytes converts a byte to a QoS type.
func QoSFromBytes(qos byte) QoS {
	switch qos {
	case 0:
		return AtMostOnce
	case 1:
		return AtLeastOnce
	default:
		return ExactlyOnce
	}
}

// ValidateQoS checks if the provided QoS value is valid.
func ValidateQoS(qos QoS) bool {
	return qos == AtMostOnce || qos == AtLeastOnce || qos == ExactlyOnce
}
