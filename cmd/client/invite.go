package main

import (
	"github.com/globalcyberalliance/domain-trust-go/v2/model"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newInvitesCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "invites",
		Short:  "Interact with invites",
		PreRun: adminCheck,
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
		},
	}

	cmd.AddCommand(newInvitesCreateCMD())
	cmd.AddCommand(newInvitesDeleteCMD())
	cmd.AddCommand(newInvitesFindCMD())
	cmd.AddCommand(newInvitesGetCMD())

	return cmd
}

func newInvitesCreateCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create invite",
		Example: "  client invites create\n  client invites create --email=dev@gcai.dev",
		Args:    cobra.ExactArgs(0),
		PreRun:  adminCheck,
		Run: func(cmd *cobra.Command, _ []string) {
			var invite model.Invite

			cmd.Flags().Visit(func(flag *pflag.Flag) {
				switch flag.Name {
				case "email":
					invite.UserEmail = flag.Value.String()
				case "firstName":
					invite.UserFirstName = flag.Value.String()
				case "lastName":
					invite.UserLastName = flag.Value.String()
				case "organizationID":
					invite.UserOrganizationID = flag.Value.String()
				case "role":
					invite.UserRole = flag.Value.String()
				}
			})

			if err := apiClient.CreateInvite(cmd.Context(), &invite); err != nil {
				log.Fatal().Err(err).Msg("Failed to create invite")
			}

			log.Info().Msg("Created invite successfully!")
		},
	}

	cmd.Flags().String("email", "", "Set invite's email")
	cmd.Flags().String("firstName", "", "Set invite's first name")
	cmd.Flags().String("lastName", "", "Set invite's last name")
	cmd.Flags().String("organizationID", "", "Set invite's organizationID")
	cmd.Flags().String("role", "", "Set invite's role")
	_ = markFlagsRequired(cmd, "email", "name")

	return cmd
}

func newInvitesDeleteCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete invite",
		Example: "  client invites delete :id\n  client invites delete 2",
		Args:    cobra.ExactArgs(1),
		PreRun:  adminCheck,
		Run: func(cmd *cobra.Command, args []string) {
			if err := apiClient.DeleteInvite(cmd.Context(), args[0]); err != nil {
				log.Fatal().Err(err).Msg("Failed to delete invite")
			}

			log.Info().Msg("Successfully deleted invite!")
		},
	}

	return cmd
}

func newInvitesFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "find",
		Short:  "Find invites",
		PreRun: adminCheck,
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.InviteFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			invites, err := apiClient.FindInvites(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find invites")
			}

			printToConsole(invites)
		},
	}

	cmd.Flags().String("organizationID", "", "Filter invites by organization ID")

	return cmd
}

func newInvitesGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get invite",
		Example: "  client invites get :id\n  client invites get 3",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			invite, err := apiClient.FindInviteByID(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get invite " + args[0])
			}

			printToConsole(invite)
		},
	}

	return cmd
}
