package client

import (
	"context"
	"fmt"

	"github.com/globalcyberalliance/domain-trust-go/v2/model"
)

func (c *Client) CreateInvite(ctx context.Context, invite *model.Invite) error {
	body, err := c.marshal(map[string]*model.Invite{"invite": invite})
	if err != nil {
		return fmt.Errorf("marshal invite: %w", err)
	}

	if _, err = c.POST(ctx, "invites", body, nil); err != nil {
		return fmt.Errorf("create invite: %w", err)
	}

	return nil
}

func (c *Client) DeleteInvite(ctx context.Context, inviteID string) error {
	if _, err := c.DELETE(ctx, "invites/"+inviteID, nil); err != nil {
		return fmt.Errorf("delete invite: %w", err)
	}

	return nil
}

func (c *Client) FindInvites(ctx context.Context, filter *model.InviteFilter) ([]*model.Invite, error) {
	query := structToQueryParams(filter)

	var response struct {
		Invites []*model.Invite `json:"invites"`
	}

	if _, err := c.GET(ctx, "invites?"+query, &response); err != nil {
		return nil, fmt.Errorf("find invites: %w", err)
	}

	return response.Invites, nil
}

func (c *Client) FindInviteByID(ctx context.Context, id string) (*model.Invite, error) {
	var response struct {
		Invite *model.Invite `json:"invite"`
	}

	if _, err := c.GET(ctx, "invites/"+id, &response); err != nil {
		return nil, fmt.Errorf("find invite: %w", err)
	}

	return response.Invite, nil
}
