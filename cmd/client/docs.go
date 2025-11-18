package main

import (
	"os/exec"
	"runtime"

	dt "github.com/globalcyberalliance/domain-trust-go/v2"
	"github.com/spf13/cobra"
)

func newDocsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Access the OpenAPI docs",
		Run: func(cmd *cobra.Command, _ []string) {
			var cmdArgs []string
			var cmdName string

			switch os := runtime.GOOS; os {
			case "darwin":
				cmdName = "open"
			case "windows":
				cmdArgs = []string{"/c", "start"}
				cmdName = "cmd"
			default: // Assume Linux or Unix.
				cmdName = "xdg-open"
			}

			cmdArgs = append(cmdArgs, dt.DocsURL)

			execCMD := exec.CommandContext(cmd.Context(), cmdName, cmdArgs...)
			if err := execCMD.Start(); err != nil {
				log.Fatal().Err(err).Msg("Failed to open browser to " + dt.DocsURL)
			}

			log.Info().Msg("Browser opened to " + dt.DocsURL)
		},
	}

	return cmd
}
