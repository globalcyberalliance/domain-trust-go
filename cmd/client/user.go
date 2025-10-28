package main

import (
	"github.com/globalcyberalliance/domain-trust-go/model"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newUserCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Show your current user information",
		Run: func(cmd *cobra.Command, _ []string) {
			user, err := apiClient.FindSessionUser(cmd.Context())
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get user info")
			}

			printToConsole(user)
		},
	}

	return cmd
}

func newUsersCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "users",
		Short:  "Interact with users",
		PreRun: adminCheck,
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
		},
	}

	cmd.AddCommand(newUsersDeleteCMD())
	cmd.AddCommand(newUsersFindCMD())
	cmd.AddCommand(newUsersGetCMD())
	cmd.AddCommand(newUsersUpdateCMD())

	return cmd
}

func newUsersDeleteCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete user",
		Example: "  client users delete :id\n  client users delete 2",
		Args:    cobra.ExactArgs(1),
		PreRun:  adminCheck,
		Run: func(cmd *cobra.Command, args []string) {
			if err := apiClient.DeleteUser(cmd.Context(), args[0]); err != nil {
				log.Fatal().Err(err).Msg("Failed to delete user")
			}

			printToConsole("User " + args[0] + " successfully deleted!")
		},
	}

	return cmd
}

func newUsersFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "find",
		Short:  "Find users",
		PreRun: adminCheck,
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.UserFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				panic(err)
			}

			users, err := apiClient.FindUsers(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find users")
			}

			if len(users) == 0 {
				log.Warn().Msg("No users found")
				return
			}

			printToConsole(users)
		},
	}

	cmd.Flags().String("email", "", "Filter users by email")
	cmd.Flags().String("firstName", "", "Filter users by first name")
	cmd.Flags().String("lastName", "", "Filter users by last name")
	cmd.Flags().String("organizationID", "", "Filter users by organization ID")
	cmd.Flags().String("role", "", "Filter users by role")

	return cmd
}

func newUsersGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get user",
		Example: "  client users get :id\n  client users get 3",
		Args:    cobra.ExactArgs(1),
		PreRun:  adminCheck,
		Run: func(cmd *cobra.Command, args []string) {
			user, err := apiClient.FindUserByID(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get user " + args[0])
			}

			printToConsole(user)
		},
	}

	return cmd
}

func newUsersUpdateCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update user",
		Example: "  client users update :id\n  client users update 3 --email=dev@gcai.dev",
		Args:    cobra.ExactArgs(1),
		PreRun:  adminCheck,
		Run: func(cmd *cobra.Command, args []string) {
			var update *model.UserUpdate

			cmd.Flags().Visit(func(flag *pflag.Flag) {
				val := flag.Value.String()
				switch flag.Name {
				case "email":
					update.Email = &val
				case "firstName":
					update.FirstName = &val
				case "lastName":
					update.LastName = &val
				case "password":
					update.Password = &val
				case "role":
					update.Role = &val
				}
			})

			user, err := apiClient.UpdateUser(cmd.Context(), args[0], update)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to update user")
			}

			printToConsole(user)
		},
	}

	cmd.Flags().String("email", "", "Update user's email")
	cmd.Flags().String("firstName", "", "Update user's first name")
	cmd.Flags().String("lastName", "", "Update user's last name")
	cmd.Flags().String("password", "", "Update user's password")
	cmd.Flags().String("role", "", "Update user's role")

	return cmd
}
