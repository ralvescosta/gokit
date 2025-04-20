package main

import (
	"github.com/spf13/cobra"

	"github.com/ralvescosta/gokit/examples/rmq_consumer/cmd"
)

// rootCmd represents the base command when called without any subcommands
var root = &cobra.Command{
	Use:     "app",
	Short:   "Observably project",
	Version: "0.0.1",
}

func main() {
	root.AddCommand(cmd.ConsumerCmd)

	root.Execute()
}
