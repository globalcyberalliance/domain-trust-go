package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/globalcyberalliance/domain-trust-go/model"
)

func (c *Client) CreateDomains(ctx context.Context, domains ...*model.DomainSubmission) ([]*model.DomainError, error) {
	body, err := c.marshal(map[string][]*model.DomainSubmission{"domains": domains})
	if err != nil {
		return nil, fmt.Errorf("marshal domains: %w", err)
	}

	var response struct {
		Errors []*model.DomainError `json:"errors"`
	}

	if _, err = c.POST(ctx, "domains", body, &response); err != nil {
		return nil, fmt.Errorf("create domains: %w", err)
	}

	return response.Errors, nil
}

func (c *Client) FindDomains(ctx context.Context, filter *model.DomainFilter) ([]*model.Domain, error) {
	query := structToQueryParams(filter)

	var response struct {
		Domains []*model.Domain `json:"domains"`
	}

	if _, err := c.GET(ctx, "domains?"+query, &response); err != nil {
		return nil, fmt.Errorf("find domains: %w", err)
	}

	return response.Domains, nil
}

func (c *Client) FindDomainsPaged(ctx context.Context, filter *model.DomainFilter) (*Iterator[*model.Domain], error) {
	fetch := func(ctx context.Context, pageToken string) ([]*model.Domain, string, error) {
		q := structToQueryParams(filter)
		if pageToken != "" {
			q += "&pageToken=" + url.QueryEscape(pageToken)
		}

		var resp struct {
			Domains       []*model.Domain `json:"domains"`
			NextPageToken string          `json:"nextPageToken"`
		}

		if _, err := c.GET(ctx, "domains?"+q, &resp); err != nil {
			return nil, "", fmt.Errorf("find domains: %w", err)
		}

		return resp.Domains, resp.NextPageToken, nil
	}

	// Initialize iterator (fetch first page lazily).
	return &Iterator[*model.Domain]{
		ctx:       ctx,
		fetchPage: fetch,
		index:     0,
		nextToken: "",
		page:      nil,
	}, nil
}
