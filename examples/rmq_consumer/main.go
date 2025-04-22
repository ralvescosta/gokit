// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package main

import (
	"github.com/spf13/cobra"

	"github.com/ralvescosta/gokit/examples/rmq_consumer/cmd"
)

// rootCmd represents the base command when called without any subcommands
var root = &cobra.Command{
	Use:     "app",
	Short:   "RabbitMQ Consumer Example",
	Version: "0.0.1",
}

func main() {
	root.AddCommand(cmd.ConsumerCmd)

	root.Execute()
}
