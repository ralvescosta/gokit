// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package mqtt

// QoS represents the Quality of Service levels for MQTT messages.
// It defines the reliability of message delivery between the client and broker.
// - AtMostOnce: Messages are delivered at most once (0).
// - AtLeastOnce: Messages are delivered at least once (1).
// - ExactlyOnce: Messages are delivered exactly once (2).
type QoS byte

const (
	// AtMostOnce represents QoS level 0.
	AtMostOnce QoS = 0
	// AtLeastOnce represents QoS level 1.
	AtLeastOnce QoS = 1
	// ExactlyOnce represents QoS level 2.
	ExactlyOnce QoS = 2
)

// LogMessage formats and returns a log message with a consistent prefix for MQTT operations.
func LogMessage(msg ...string) string {
	prefix := "[gokit::mqtt] "
	for _, s := range msg {
		prefix += s
	}
	return prefix
}

// QoSFromBytes converts a byte value to a QoS type.
// It maps the byte to one of the defined QoS levels.
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
// It ensures the QoS is one of the defined levels: AtMostOnce, AtLeastOnce, or ExactlyOnce.
func ValidateQoS(qos QoS) bool {
	return qos == AtMostOnce || qos == AtLeastOnce || qos == ExactlyOnce
}
