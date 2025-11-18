package main

import (
	"time"

	"github.com/globalcyberalliance/domain-trust-go/v2/model"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newAPIKeysCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apiKeys",
		Aliases: []string{"apikeys"},
		Short:   "Interact with api keys",
	}

	cmd.AddCommand(newAPIKeysCreateCMD())
	cmd.AddCommand(newAPIKeysDeleteCMD())
	cmd.AddCommand(newAPIKeysFindCMD())
	cmd.AddCommand(newAPIKeysGetCMD())

	return cmd
}

func newAPIKeysCreateCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create api key",
		Example: "  client apikeys create\n  client apikeys create --email=dev@gcai.dev",
		Args:    cobra.ExactArgs(0),
		PreRun:  adminCheck,
		Run: func(cmd *cobra.Command, _ []string) {
			var apiKey *model.APIKey

			cmd.Flags().Visit(func(flag *pflag.Flag) {
				switch flag.Name {
				case "expiry":
					expiry, err := cast.StringToDate(flag.Value.String())
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse expiry date")
					}

					if expiry.Before(time.Now()) {
						log.Fatal().Msg("Expiry can't be in the past")
					}

					apiKey.Expiry = expiry
				case "userID":
					if cfg.UserRole != model.UserRoleAdmin {
						log.Fatal().Msg("Only admins can create api keys for other users")
					}

					apiKey.UserID = flag.Value.String()
				}
			})

			if err := apiClient.CreateAPIKey(cmd.Context(), apiKey); err != nil {
				log.Fatal().Err(err).Msg("Failed to create api key")
			}

			printToConsole(apiKey)
		},
	}

	cmd.Flags().String("expiry", "", "Set expiry date (e.g. "+time.Now().Format("2006-01-02")+")")
	cmd.Flags().Int64("userID", 0, "Set user id")
	_ = markFlagsRequired(cmd, "expiry")

	return cmd
}

func newAPIKeysDeleteCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete api key",
		Example: "  client apikeys delete :id\n  client apikeys delete c318fc9b-d343-4a8f-b3cb-1a18773ddf05",
		Args:    cobra.ExactArgs(1),
		PreRun:  adminCheck,
		Run: func(cmd *cobra.Command, args []string) {
			if err := apiClient.DeleteAPIKey(cmd.Context(), args[0]); err != nil {
				log.Fatal().Err(err).Msg("Failed to delete api key")
			}

			log.Info().Msg("Successfully deleted API key!")
		},
	}

	return cmd
}

func newAPIKeysFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find api keys",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.APIKeyFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				panic(err)
			}

			apiKeys, err := apiClient.FindAPIKeys(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find api keys")
			}

			if len(apiKeys) == 0 {
				log.Warn().Msg("No api keys found")
				return
			}

			printToConsole(apiKeys)
		},
	}

	// Environment.
	cmd.Flags().String("environment", "", "Which environment your requests should be made against (production|test)")
	cmd.Flags().String("expiryAfter", "", "Filter for API keys which expire after the given timestamp (e.g. 2022-08-30T00:00:00.001Z)")
	cmd.Flags().String("expiryBefore", "", "Filter for API keys which expire before the given timestamp (e.g. 2022-08-30T00:00:00.001Z)")
	cmd.Flags().String("userID", "", "Filter for API keys belonging to a specific user")

	return cmd
}

func newAPIKeysGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get api key",
		Example: "  client apiKeys get :id\n  client apiKeys get 19ad4d0e-569d-4ee7-9ee5-c24594a1acd9",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			apiKey, err := apiClient.FindAPIKeyByID(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get apiKey " + args[0])
			}

			printToConsole(apiKey)
		},
	}

	return cmd
}
