package rabbitmq

func LogMessage(msg ...string) string {
	f := "[gokit::rabbitmq] "

	for _, s := range msg {
		f += s
	}

	return f
}
