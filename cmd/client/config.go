package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newConfigCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure your DSS instance",
	}

	cmd.AddCommand(newConfigGetCMD())
	cmd.AddCommand(newConfigSetCMD())
	cmd.AddCommand(newConfigShowCMD())

	return cmd
}

func newConfigGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get a config value",
		Example: "  client config get apikey",
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			switch strings.ToLower(args[0]) {
			case "apikey":
				printToConsole("api key: " + cfg.APIKey)
			case "useremail":
				printToConsole("user email: " + cfg.UserEmail)
			case "userpass":
				printToConsole("user pass: " + cfg.UserPass)
			case "userrole":
				printToConsole("user role: " + cfg.UserRole)
			default:
				log.Fatal().Msg("unknown config key")
			}
		},
	}

	return cmd
}

func newConfigSetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Short:   "Set a config value",
		Example: "  client config set apikey 019a0dd4-11a5-7477-91a8-538b1bc334e4",
		Args:    cobra.ExactArgs(2), //nolint:mnd // Unnecessary.
		Run: func(_ *cobra.Command, args []string) {
			switch strings.ToLower(args[0]) {
			case "apikey":
				cfg.APIKey = args[1]
			case "useremail":
				cfg.UserEmail = args[1]
			case "userpass":
				cfg.UserPass = args[1]
			case "userrole":
				cfg.UserRole = args[1]
			default:
				log.Fatal().Msg("unknown config key")
			}

			if err := cfg.Save(); err != nil {
				log.Fatal().Err(err).Msg("unable to save config")
			}

			log.Info().Msg("config updated")
		},
	}

	return cmd
}

func newConfigShowCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Print full config",
		Example: "  dss config show",
		Args:    cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			printToConsole(cfg)
		},
	}

	return cmd
}

type Config struct {
	dir       string
	path      string
	APIKey    string `json:"apiKey" yaml:"apiKey"`
	UserEmail string `json:"userEmail" yaml:"userEmail"`
	UserPass  string `json:"userPass,omitempty" yaml:"userPass,omitempty"`
	UserRole  string `json:"userRole" yaml:"userRole"`
}

func newConfig(directory string) (*Config, error) {
	config := Config{
		dir:       directory,
		path:      directory + slash + "config.yml",
		APIKey:    "",
		UserEmail: "",
		UserPass:  "",
		UserRole:  "",
	}

	if err := config.Load(); err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return &config, nil
}

func (c *Config) Load() error {
	// Create config if it doesn't exist.
	if _, err := os.Stat(c.path); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(c.dir, 0o750); err != nil {
			log.Fatal().Err(err).Msg("Failed to create config directory")
		}

		if err = c.Save(); err != nil {
			return err
		}
	}

	// Read config.
	configData, err := os.ReadFile(c.path)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to read config file")
	}

	if err = yaml.Unmarshal(configData, &c); err != nil {
		log.Fatal().Err(err).Msg("unable to unmarshal config values")
	}

	return nil
}

func (c *Config) Save() error {
	configData, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal config values: %w", err)
	}

	if err = os.WriteFile(c.path, configData, 0o600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
