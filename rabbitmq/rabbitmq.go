package rabbitmq

func LogMessage(msg string) string {
	return "[gokit::rabbitmq] " + msg
}
