package mqtt

type QoS byte

var (
	AtMostOnce  QoS = 0
	AtLeastOnce QoS = 1
	ExactlyOnce QoS = 2
)

func LogMessage(msg ...string) string {
	f := "[gokit::mqtt] "

	for _, s := range msg {
		f += s
	}

	return f
}

func QoSFromBytes(qos byte) QoS {
	if qos == byte(0) {
		return AtMostOnce
	}

	if qos == byte(1) {
		return AtLeastOnce
	}

	return ExactlyOnce
}

func ValidateQoS(qos QoS) bool {
	if qos == AtMostOnce || qos == AtLeastOnce || qos == ExactlyOnce {
		return true
	}

	return false
}
