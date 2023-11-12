package mqtt

import (
	myQTT "github.com/eclipse/paho.mqtt.golang"
)

func New() {
	opts := myQTT.NewClientOptions()

	opts.AddBroker("")
	opts.SetClientID("")
	opts.SetUsername("")
	opts.SetPassword("")

	client := myQTT.NewClient(opts)

	client.Connect()

	println(opts)
}
