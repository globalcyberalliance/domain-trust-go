package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	dt "github.com/globalcyberalliance/domain-trust-go/v2"
	"github.com/globalcyberalliance/domain-trust-go/v2/model"
	"github.com/rs/zerolog"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	defaultTimeout = 60 * time.Second
)

var (
	apiClient                     *dt.Client
	cfg                           *Config
	log                           zerolog.Logger
	debug, prettyLog, writeToFile bool
	timeout                       time.Duration
	format, logLevel              string
	limit                         uint64
	slash                         = string(os.PathSeparator)
)

func main() {
	rootCMD := newRootCMD()
	rootCMD.AddCommand(newAPIKeysCMD())
	rootCMD.AddCommand(newConfigCMD())
	rootCMD.AddCommand(newDocsCMD())
	rootCMD.AddCommand(newDomainsCMD())
	rootCMD.AddCommand(newInvitesCMD())
	rootCMD.AddCommand(newLoginCMD())
	rootCMD.AddCommand(newUserCMD())
	rootCMD.AddCommand(newUsersCMD())
	rootCMD.AddCommand(newVersionCMD())

	if err := rootCMD.Execute(); err != nil {
		panic(err)
	}
}

func newRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "domain-trust Client",
		Long:  `Interact with the domain-trust API`,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			logger, err := newLogger(logLevel)
			if err != nil {
				fmt.Printf("Unable to setup logger: %s: %v\n", logLevel, err)
				os.Exit(1)
			}

			log = logger

			configDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal().Err(err).Msg("unable to retrieve user's home directory")
			}

			cfg, err = newConfig(fmt.Sprintf("%s%s.config%sdomain-trust-client", strings.TrimSuffix(configDir, slash), slash, slash))
			if err != nil {
				log.Fatal().Err(err).Msg("unable to initialize config")
			}

			apiClient = dt.New(cfg.APIKey, dt.WithDebug(debug), dt.WithTimeout(defaultTimeout))
		},
		Version: dt.Version,
	}

	cmd.PersistentFlags().String("createdAfter", "", "Only return results created after this date")
	cmd.PersistentFlags().String("createdBefore", "", "Only return results created before this date")
	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable console debugging")
	cmd.PersistentFlags().StringVarP(&format, "format", "f", "yaml", "Set the output format for CLI commands (json, jsonp, yaml)")
	cmd.PersistentFlags().Uint64VarP(&limit, "limit", "l", 0, "Limit the quantity of returned results")
	cmd.PersistentFlags().StringVar(&logLevel, "logLevel", "info", "Set log level (debug, info, warn, error, fatal, panic)")
	cmd.PersistentFlags().BoolVar(&prettyLog, "prettyLog", true, "Pretty print logs to console")
	cmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", defaultTimeout, "Specify the API HTTP timeout")
	cmd.PersistentFlags().BoolVarP(&writeToFile, "writetofile", "w", false, "Write the output to a file")

	return cmd
}

func adminCheck(_ *cobra.Command, _ []string) {
	if cfg != nil && cfg.UserRole != "" && cfg.UserRole != model.UserRoleAdmin {
		log.Fatal().Msg("Only admins may do this")
	}
}

func markFlagsRequired(cmd *cobra.Command, flags ...string) error {
	for _, flag := range flags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			return fmt.Errorf("mark required flag %s", flag)
		}
	}

	return nil
}

func marshal(data any) []byte {
	var err error
	var output []byte

	switch strings.ToLower(format) {
	case "cbor":
		output, err = cbor.Marshal(data)
	case "json":
		output, err = json.Marshal(data)
	case "jsonp":
		output, err = json.MarshalIndent(data, "", "\t")
	default:
		output, err = yaml.Marshal(data)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("unable to marshal data")
	}

	return output
}

func newLogger(logLevel string) (zerolog.Logger, error) {
	var logger zerolog.Logger

	logLevelParsed, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		return logger, fmt.Errorf("invalid log level %s", logLevel)
	}

	var logWriter io.Writer

	if prettyLog {
		logWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		logWriter = os.Stdout
	}

	logger = zerolog.New(logWriter).With().Timestamp().Logger().Level(logLevelParsed)

	return logger, nil
}

func printToConsole(data any) {
	if writeToFile {
		extension := format
		if extension == "jsonp" {
			extension = "json"
		}

		filename := cast.ToString(time.Now().Unix()) + "." + extension

		if err := printToFile(data, filename); err != nil {
			log.Fatal().Err(err).Msg("Failed to write output to file")
		}

		log.Info().Msg("Output written to " + filename)

		return
	}

	marshalledData := marshal(data)

	print(string(marshalledData)) //nolint:forbidigo
}

func printToFile(data any, file string) error {
	marshalledData := marshal(data)

	outputFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer outputFile.Close()

	if _, err = outputFile.Write(marshalledData); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
