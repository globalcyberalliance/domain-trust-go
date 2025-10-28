package main

import (
	dt "github.com/globalcyberalliance/domain-trust-go"
	"github.com/spf13/cobra"
)

func newLoginCMD() *cobra.Command {
	var email, password string

	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to the API using an existing API key, or a user's email and password.\nFor convenience, your email and password are stored locally (in plaintext) for future authentication requests.",
		Example: "  client login 2fbf890c-aa01-4f5b-907f-a9c2ba634ec2",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if email == "" || password == "" {
					email = cfg.UserEmail
					password = cfg.UserPass

					if email == "" || password == "" {
						log.Fatal().Msg("No key, email or password provided")
					}
				}
			}

			if email == "" {
				apiKey := args[0]

				apiClient = dt.New(apiKey)

				user, err := apiClient.FindSessionUser(cmd.Context())
				if err != nil {
					log.Fatal().Err(err).Msg("unable to retrieve user information")
				}

				cfg.APIKey = apiKey
				cfg.UserRole = user.Role

				if err = cfg.Save(); err != nil {
					log.Fatal().Err(err).Msg("could not save config")
				}

				printToConsole("Successfully set api key as " + apiKey)
			} else {
				apiClient = dt.New(cfg.APIKey)

				apiKey, err := apiClient.Login(cmd.Context(), email, password)
				if err != nil {
					log.Fatal().Err(err).Msg("unable to login user")
				}

				apiClient.SetAPIKey(apiKey.Key)

				user, err := apiClient.FindSessionUser(cmd.Context())
				if err != nil {
					log.Fatal().Err(err).Msg("unable to retrieve user information")
				}

				cfg.APIKey = apiKey.Key
				cfg.UserEmail = email
				cfg.UserPass = password
				cfg.UserRole = user.Role

				if err = cfg.Save(); err != nil {
					log.Fatal().Err(err).Msg("could not save config")
				}

				printToConsole("Successfully logged in!")
				printToConsole("API key set!")
			}
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Login with your email")
	cmd.Flags().StringVar(&password, "password", "", "Login with your password")

	return cmd
}
