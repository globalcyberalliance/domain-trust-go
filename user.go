package client

import (
	"context"
	"fmt"

	"github.com/globalcyberalliance/domain-trust-go/v2/model"
)

func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	if _, err := c.DELETE(ctx, "users/"+userID, nil); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func (c *Client) FindSessionUser(ctx context.Context) (*model.User, error) {
	var response struct {
		User *model.User `json:"user"`
	}

	if _, err := c.GET(ctx, "user", &response); err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return response.User, nil
}

func (c *Client) FindUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, error) {
	query := structToQueryParams(filter)

	var response struct {
		Users []*model.User `json:"users"`
	}

	if _, err := c.GET(ctx, "users?"+query, &response); err != nil {
		return nil, fmt.Errorf("find users: %w", err)
	}

	return response.Users, nil
}

func (c *Client) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	var response struct {
		User *model.User `json:"user"`
	}

	if _, err := c.GET(ctx, "users/"+id, &response); err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return response.User, nil
}

func (c *Client) UpdateUser(ctx context.Context, id string, update *model.UserUpdate) (*model.User, error) {
	body, err := c.marshal(map[string]*model.UserUpdate{"user": update})
	if err != nil {
		return nil, fmt.Errorf("marshal update: %w", err)
	}

	var response struct {
		User *model.User `json:"user"`
	}

	if _, err = c.PATCH(ctx, "users/"+id, body, &response); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return response.User, nil
}
