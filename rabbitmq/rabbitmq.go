// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

func LogMessage(msg ...string) string {
	f := "[gokit::rabbitmq] "

	for _, s := range msg {
		f += s
	}

	return f
}
