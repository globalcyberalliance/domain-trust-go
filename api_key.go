package client

import (
	"context"
	"fmt"

	"github.com/globalcyberalliance/domain-trust-go/model"
)

func (c *Client) CreateAPIKey(ctx context.Context, apiKey *model.APIKey) error {
	body, err := c.marshal(map[string]*model.APIKey{"key": apiKey})
	if err != nil {
		return fmt.Errorf("marshal api key: %w", err)
	}

	if _, err = c.POST(ctx, "keys", body, nil); err != nil {
		return fmt.Errorf("create api key: %w", err)
	}

	return nil
}

func (c *Client) DeleteAPIKey(ctx context.Context, apiKeyID string) error {
	if _, err := c.DELETE(ctx, "keys/"+apiKeyID, nil); err != nil {
		return fmt.Errorf("delete api key: %w", err)
	}

	return nil
}

func (c *Client) FindAPIKeys(ctx context.Context, filter *model.APIKeyFilter) ([]*model.APIKey, error) {
	query := structToQueryParams(filter)

	var response struct {
		APIKeys []*model.APIKey `json:"keys"`
	}

	if _, err := c.GET(ctx, "keys?"+query, &response); err != nil {
		return nil, fmt.Errorf("find api keys: %w", err)
	}

	return response.APIKeys, nil
}

func (c *Client) FindAPIKeyByID(ctx context.Context, id string) (*model.APIKey, error) {
	var response struct {
		APIKey *model.APIKey `json:"key"`
	}

	if _, err := c.GET(ctx, "keys/"+id, &response); err != nil {
		return nil, fmt.Errorf("find api key: %w", err)
	}

	return response.APIKey, nil
}
