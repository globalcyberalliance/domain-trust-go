// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type (
	DashboardMetrics struct {
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
		DomainsByProvider30d     []KV `json:"domainsByProvider30d"`
		SubmissionsByProvider    []KV `json:"submissionsByProvider"`
		SubmissionsByProvider30d []KV `json:"submissionsByProvider30d"`

		DailyDomains30d     []TimeCount `json:"dailyDomains30d"`
		DailySubmissions30d []TimeCount `json:"dailySubmissions30d"`
	}

	DashboardMetricsPublic struct {
		UniqueDomains uint64 `json:"uniqueDomains"`
		Submissions   uint64 `json:"submissions"`
		TotalPartners uint64 `json:"totalPartners"`

		DailyDomains30d     []TimeCount `json:"dailyDomains30d"`
		DailySubmissions30d []TimeCount `json:"dailySubmissions30d"`
	}

	KV struct {
		Key   string `json:"key"`
		Value uint64 `json:"value"`
	}

	RatingClassificationCount struct {
		Rating         string `json:"rating"`
		Classification string `json:"classification"`
		Count          uint64 `json:"count"`
	}

	TimeCount struct {
		Date  time.Time `json:"date"`
		Count uint64    `json:"count"`
	}
)
