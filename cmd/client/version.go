package main

import (
	"github.com/spf13/cobra"
)

func newVersionCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get API version information",
		Run: func(cmd *cobra.Command, _ []string) {
			apiVersion, err := apiClient.FindVersion(cmd.Context())
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find version info")
			}

			printToConsole(apiVersion)
		},
	}

	return cmd
}
