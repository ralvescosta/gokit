package mqtt

func LogMessage(msg ...string) string {
	f := "[gokit::mqtt] "

	for _, s := range msg {
		f += s
	}

	return f
}
