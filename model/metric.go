// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type DashboardMetrics struct {
	UniqueDomains     uint64  `json:"uniqueDomains"`
	Submissions       uint64  `json:"submissions"`
	TotalPartners     uint64  `json:"totalPartners"`
	ActivePartners30d uint64  `json:"activePartners30d"`
	PctBlocked        float64 `json:"pctBlocked"`

	DomainsByPartnerRating        []KV                        `json:"domainsByPartnerRating"`
	DomainsByClassification       []KV                        `json:"domainsByClassification"`
	DomainsByRatingClassification []RatingClassificationCount `json:"domainsByRatingClassification"`
	DomainsByActivity             []KV                        `json:"domainsByActivity"`

	DomainsByProvider        []KV `json:"domainsByProvider"`
	DomainsByProviderWeekly  []KV `json:"domainsByProviderWeekly"`
	SubmissionsByProvider    []KV `json:"submissionsByProvider"`
	SubmissionsByProvider30d []KV `json:"submissionsByProvider30d"`

	WeeklyDomains     []TimeCount        `json:"weeklyDomains"`
	WeeklySubmissions []TimeCount        `json:"weeklySubmissions"`
	OrgMetrics        []*WeeklyOrgMetric `json:"orgMetrics,omitempty"`

	ClientPlatforms []KV `json:"clientPlatforms"`
}

type DashboardMetricsPublic struct {
	UniqueDomains uint64 `json:"uniqueDomains"`
	Submissions   uint64 `json:"submissions"`
	TotalPartners uint64 `json:"totalPartners"`

	WeeklyDomains     []TimeCount `json:"weeklyDomains"`
	WeeklySubmissions []TimeCount `json:"weeklySubmissions"`
}

type KV struct {
	Key   string `json:"key"`
	Value uint64 `json:"value"`
}

type RatingClassificationCount struct {
	Rating         string `json:"rating"`
	Classification string `json:"classification"`
	Count          uint64 `json:"count"`
}

type TimeCount struct {
	Date  time.Time `json:"date"`
	Count uint64    `json:"count"`
}

type OrgUsage struct {
	OrganizationID   string `json:"organizationID"`
	OrganizationName string `json:"organizationName"`
	Consumption      int    `json:"consumption"`
	Production       int    `json:"production"`
}

type WeeklyOrgMetric struct {
	Week   time.Time  `json:"week"`
	Usages []OrgUsage `json:"usages"`
}
