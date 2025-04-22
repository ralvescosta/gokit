// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package cmd

import (
	"github.com/spf13/cobra"
)

var HTTPServerCmd = &cobra.Command{
	Use:   "http",
	Short: "HTTP Server Command",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctn, err := NewContainer()
		if err != nil {
			return err
		}

		ctn.Logger.Debug("Starting HTTP Server...")

		ctn.BooksHandlers.Install(ctn.HTTPServer)

		return ctn.HTTPServer.Run()
	},
}
