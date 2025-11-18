package client

import (
	"context"
	"fmt"

	"github.com/globalcyberalliance/domain-trust-go/v2/model"
)

func (c *Client) Login(ctx context.Context, email string, password string) (*model.APIKey, error) {
	body, err := c.marshal(map[string]*model.Login{"login": {Email: email, Password: password}})
	if err != nil {
		return nil, fmt.Errorf("marshal login request: %w", err)
	}

	var response struct {
		Key *model.APIKey `json:"key"`
	}

	if _, err = c.POST(ctx, "auth/login", body, &response); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return response.Key, nil
}
