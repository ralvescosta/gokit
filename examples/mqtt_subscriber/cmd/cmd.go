// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package cmd

import (
	"github.com/spf13/cobra"
)

var ConsumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consumer Command",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctn, err := NewContainer()
		if err != nil {
			return err
		}

		ctn.Logger.Debug("Starting MQTT Consumer...")

		ctn.BasicConsumer.Install(ctn.Dispatcher)
		ctn.Dispatcher.ConsumeBlocking(ctn.Sig)

		return nil
	},
}
