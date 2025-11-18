package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/globalcyberalliance/domain-trust-go/v2/model"
	"github.com/spf13/cobra"
)

func newDomainsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domains",
		Short: "Interact with domains",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
		},
	}

	cmd.AddCommand(newDomainsCreateCMD())
	cmd.AddCommand(newDomainsFindCMD())

	return cmd
}

func newDomainsCreateCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create domains",
		Example: "  client domains create\n  client users create < input.csv",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			// Read CSV from stdin.
			domains, err := readDomainsCSV(os.Stdin)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to parse CSV from stdin")
			}
			if len(domains) == 0 {
				log.Fatal().Msg("No domain rows found in CSV")
			}

			errs, err := apiClient.CreateDomains(cmd.Context(), domains...)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create domains")
			}

			if len(errs) > 0 {
				log.Info().Msg("Domains created with errors")
				for _, domain := range errs {
					printToConsole(domain)
				}

				return
			}

			log.Info().Msg("Domains created!")
		},
	}

	return cmd
}

func newDomainsFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find domains",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.DomainFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			findAll, err := cmd.Flags().GetBool("all")
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get flag 'all'")
			}

			var domains []*model.Domain

			// If findAll is false, do a simple lookup and return the results.
			if !findAll {
				domains, err = apiClient.FindDomains(cmd.Context(), &filter)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to find domains")
				}

				if len(domains) == 0 {
					log.Warn().Msg("No domains found")
					return
				}

				printToConsole(domains)
				return
			}

			go func() {
				ticker := time.NewTicker(10 * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-cmd.Context().Done():
						return
					case <-ticker.C:
						log.Info().
							Int("domainsFound", len(domains)).
							Msg("Tracking domains")
					}
				}
			}()

			log.Info().Msg("Starting domain retrieval...")

			// Paginate over results.
			filter.MetadataFilter.Limit = model.MaxMetadataLimit

			domainIterator, err := apiClient.FindDomainsPaged(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find all domains")
			}

			domains = make([]*model.Domain, 0, model.MaxMetadataLimit)

			for domainIterator.Next() {
				domains = append(domains, domainIterator.Value())
			}

			if domainIterator.Err() != nil {
				log.Fatal().Err(domainIterator.Err()).Msg("Failed to page domains")
			}

			if len(domains) == 0 {
				log.Warn().Msg("No domains found")
				return
			}

			log.Info().Int("domainsFound", len(domains)).Msg("Successfully retrieved domains")

			printToConsole(domains)
		},
	}

	cmd.Flags().Bool("all", false, "Automatically paginate through the results")
	cmd.Flags().String("organizationID", "", "A unique identifier for the organization")

	// Domain details.
	cmd.Flags().String("domain", "", "The fully qualified domain name (FQDN) of the submission")
	cmd.Flags().String("sld", "", "Second-level domain portion of the FQDN")
	cmd.Flags().String("tld", "", "Top-level domain portion of the FQDN")
	cmd.Flags().String("rootDomain", "", "The root domain extracted from the FQDN")
	cmd.Flags().String("subDomain", "", "The subdomain portion of the domain, if applicable")

	// Organization details.
	cmd.Flags().String("providerName", "", "The name of the provider organization")
	cmd.Flags().String("providerRating", "", "The rating of the provider organization (trial|predictive|low-confidence|med-confidence|high-confidence)")
	cmd.Flags().String("providerRatingAbove", "", "Filter for providers with a rating above the given rating (trial|predictive|low-confidence|med-confidence|high-confidence)")
	cmd.Flags().String("providerRole", "", "The role/type of the provider organization (registrar|registry|reseller|other|ICANN)")

	// Submission details.
	cmd.Flags().String("abuseType", "", "Type of abuse associated with the domain (botnets|malware|pharming|phishing|spam)")
	cmd.Flags().String("activity", "", "The activity for the submission (active|suspended|non-existent|taken-down|blocked)")
	cmd.Flags().String("classification", "", "The classification of the submission (definitely-malicious|probably-malicious|possibly-malicious|definitely-clean)")
	cmd.Flags().Bool("isBlocked", false, "Whether the FQDN is blocked by Quad9")
	cmd.Flags().String("reportType", "", "The type of the submission (brand-spoof|fraud)")
	cmd.Flags().String("source", "", "The source of the submission (self-reported|external-reported)")
	cmd.Flags().String("urls", "", "Reported URLs for this domain")

	return cmd
}

// readDomainsCSV parses a CSV with a header row into domainDTOs.
// Recognized headers (case-insensitive):
// abuseType, activity, classification, comments, dateIdentified, domain,
// reportType, rootDomain, source, sourceName, urls.
func readDomainsCSV(r io.Reader) ([]*model.DomainSubmission, error) {
	cr := csv.NewReader(bufio.NewReader(r))
	cr.TrimLeadingSpace = true
	cr.FieldsPerRecord = -1 // Allow variable columns.

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}

	// Map header to index.
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	var out []*model.DomainSubmission
	for {
		rec, rErr := cr.Read()
		if rErr != nil {
			if errors.Is(rErr, io.EOF) {
				break
			}
			return nil, fmt.Errorf("read csv record: %w", rErr)
		}

		// Helper to safely fetch a field by name.
		get := func(name string) string {
			if j, ok := idx[name]; ok && j < len(rec) {
				return strings.TrimSpace(rec[j])
			}
			return ""
		}

		// Split urls by comma/semicolon/whitespace.
		urlsRaw := get("urls")
		var urls []string
		if urlsRaw != "" {
			for _, u := range splitMulti(urlsRaw, []rune{',', ';', ' ', '\t', '\n'}) {
				if u != "" {
					urls = append(urls, u)
				}
			}
		}

		d := model.DomainSubmission{
			AbuseType:      get("abusetype"),
			Activity:       get("activity"),
			Classification: get("classification"),
			Comments:       get("comments"),
			Domain:         get("domain"),
			ReportType:     get("reporttype"),
			Source:         get("source"),
			SourceName:     get("sourcename"),
			URLs:           urls,
		}

		// Require a domain value.
		if d.Domain == "" {
			// Skip empty rows silently.
			continue
		}

		out = append(out, &d)
	}

	return out, nil
}

// splitMulti splits s by any of the given runes.
func splitMulti(s string, seps []rune) []string {
	if s == "" {
		return nil
	}

	fields := strings.FieldsFunc(s, func(r rune) bool {
		for _, sep := range seps {
			if r == sep {
				return true
			}
		}
		return false
	})

	trimmed := make([]string, 0, len(fields))
	for _, v := range fields {
		v = strings.TrimSpace(v)
		if v != "" {
			trimmed = append(trimmed, v)
		}
	}

	return trimmed
}
