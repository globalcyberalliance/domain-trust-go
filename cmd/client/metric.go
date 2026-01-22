package main

import (
	"github.com/spf13/cobra"
)

func newMetricsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "metrics",
		Short:  "Interact with metrics",
		PreRun: adminCheck,
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
		},
	}

	cmd.AddCommand(newMetricsFindCMD())

	return cmd
}

func newMetricsFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "find",
		Short:  "Find metrics",
		PreRun: adminCheck,
		Run: func(cmd *cobra.Command, _ []string) {
			public, err := cmd.Flags().GetBool("public")
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get flag 'public'")
			}

			if public {
				metrics, err := apiClient.FindMetricsPublic(cmd.Context())
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to find metrics (public)")
				}

				printToConsole(metrics)
				return
			}

			metrics, err := apiClient.FindMetrics(cmd.Context())
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find metrics")
			}

			printToConsole(metrics)
		},
	}

	cmd.Flags().Bool("public", false, "Retrieve public metrics")

	return cmd
}
