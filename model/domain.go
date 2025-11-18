// Code generated from internal models; DO NOT EDIT.
package model

import "time"

const (
	DomainAbuseTypeBotnets  = "botnets"
	DomainAbuseTypeMalware  = "malware"
	DomainAbuseTypePharming = "pharming"
	DomainAbuseTypePhishing = "phishing"
	DomainAbuseTypeSpam     = "spam"

	DomainActivityActive      = "active"       // Active domain.
	DomainActivitySuspended   = "suspended"    // Suspended domain.
	DomainActivityNonExistent = "non-existent" // Domain doesn't exist.
	DomainActivityTakenDown   = "taken-down"   // Taken down.
	DomainActivityBlocked     = "blocked"      // Blocked by Quad9.

	DomainClassificationDefinitelyMalicious = "definitely-malicious" // In the provider's opinion.
	DomainClassificationProbablyMalicious   = "probably-malicious"   // In the provider's opinion.
	DomainClassificationPossiblyMalicious   = "possibly-malicious"   // In the provider's opinion.
	DomainClassificationDefinitelyClean     = "definitely-clean"     // Intended for whitelists, negative false positives, or investigations underway.

	DomainReportTypeBrandSpoof = "brand-spoof"
	DomainReportTypeFraud      = "fraud"

	DomainSourceExternal = "external-reported"
	DomainSourceInternal = "self-reported"
)

type (
	Domain struct {
		DomainSubmission

		OtherProviders []*DomainSubmission `json:"otherProviders,omitempty" yaml:"otherProviders,omitempty"`
		Whitelist      []*DomainWhitelist  `json:"whitelist,omitempty" yaml:"whitelist,omitempty"`
	}

	DomainError struct {
		Domain string `json:"domain"`
		Error  string `json:"error"`
	}

	DomainFilter struct {
		MetadataFilter

		OrganizationID string `json:"organizationID,omitempty" query:"organizationID"`

		CreatedAfter  time.Time `json:"createdAfter,omitzero" query:"createdAfter"`
		CreatedBefore time.Time `json:"createdBefore,omitzero" query:"createdBefore"`

		Domain string `json:"domain,omitempty" query:"domain"`

		SLD                    string    `json:"sld,omitempty" query:"sld"`
		TLD                    string    `json:"tld,omitempty" query:"tld"`
		RootDomain             string    `json:"rootDomain,omitempty" query:"rootDomain"`
		Subdomain              string    `json:"subDomain,omitempty" query:"subDomain"`
		RegistrationDateAfter  time.Time `json:"registrationDateAfter,omitzero" query:"registrationDateAfter"`
		RegistrationDateBefore time.Time `json:"registrationDateBefore,omitzero" query:"registrationDateBefore"`

		ProviderName        string `json:"providerName,omitempty" query:"providerName"`
		ProviderRating      string `json:"providerRating,omitempty" query:"providerRating"`
		ProviderRatingAbove string `json:"providerRatingAbove,omitempty" query:"providerRatingAbove"`
		ProviderRole        string `json:"providerRole,omitempty" query:"providerRole"`

		AbuseType            string    `json:"abuseType,omitempty" query:"abuseType"`
		Activity             string    `json:"activity,omitempty" query:"activity"`
		Classification       string    `json:"classification,omitempty" query:"classification"`
		DateIdentifiedAfter  time.Time `json:"dateIdentifiedAfter,omitzero" query:"dateIdentifiedAfter"`
		DateIdentifiedBefore time.Time `json:"dateIdentifiedBefore,omitzero"doc:"The query:"dateIdentifiedBefore"`
		OnlyBlocked          bool      `json:"onlyBlocked,omitempty" query:"onlyBlocked"`
		OnlyUnblocked        bool      `json:"onlyUnblocked,omitempty" query:"onlyUnblocked"`
		ReportType           string    `json:"reportType,omitempty" query:"reportType"`
		Source               string    `json:"source,omitempty" query:"source"`
		URLs                 string    `json:"urls,omitempty" query:"urls"`
	}

	DomainUpdate struct {
		OrganizationID string  `json:"organizationID"`
		ProviderName   *string `json:"providerName"`
		ProviderRating *string `json:"providerRating"`
		ProviderRole   *string `json:"providerRole"`
	}

	DomainSubmission struct {
		Created        time.Time `json:"created,omitzero" yaml:"created,omitempty"`
		ID             string    `json:"id,omitempty" yaml:"id,omitempty"`
		OrganizationID string    `json:"providerID,omitempty" yaml:"providerID,omitempty"`

		Domain           string    `json:"domain" yaml:"domain"`
		SLD              string    `json:"sld,omitempty" yaml:"sld,omitempty"`
		TLD              string    `json:"tld,omitempty" yaml:"tld,omitempty"`
		RootDomain       string    `json:"rootDomain,omitempty" yaml:"rootDomain,omitempty"`
		Subdomain        string    `json:"subDomain,omitempty" yaml:"subDomain,omitempty"`
		RegistrationDate time.Time `json:"registrationDate,omitzero" yaml:"registrationDate,omitempty"`

		ProviderName   string `json:"providerName,omitempty" yaml:"providerName,omitempty"`
		ProviderRating string `json:"providerRating,omitempty" yaml:"providerRating,omitempty"`
		ProviderRole   string `json:"providerRole,omitempty" yaml:"providerRole,omitempty"`

		AbuseType      string    `json:"abuseType,omitempty" yaml:"abuseType,omitempty"`
		Activity       string    `json:"activity" yaml:"activity"`
		Classification string    `json:"classification" yaml:"classification"`
		Comments       string    `json:"comments,omitempty" yaml:"comments,omitempty"`
		DateIdentified time.Time `json:"dateIdentified,omitzero" yaml:"dateIdentified,omitempty"`
		IsBlocked      bool      `json:"isBlocked" yaml:"isBlocked"`
		ReportType     string    `json:"reportType,omitempty" yaml:"reportType,omitempty"`
		Source         string    `json:"source" yaml:"source"`
		SourceName     string    `json:"sourceName,omitempty" yaml:"sourceName,omitempty"`
		URLs           []string  `json:"urls,omitzero" yaml:"urls,omitempty"`
	}

	DomainWhitelist struct {
		Created time.Time `json:"created" yaml:"created"`

		ProviderName string `json:"providerName" yaml:"providerName"`
		SourceName   string `json:"sourceName,omitempty" yaml:"sourceName,omitempty"`
	}
)
