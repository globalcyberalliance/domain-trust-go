package client

import (
	"context"
	"fmt"

	"github.com/globalcyberalliance/domain-trust-go/v2/model"
)

func (c *Client) FindMetrics(ctx context.Context) (*model.DashboardMetrics, error) {
	var response struct {
		Metrics *model.DashboardMetrics `json:"metrics"`
	}

	if _, err := c.GET(ctx, "metrics", &response); err != nil {
		return nil, fmt.Errorf("find metrics: %w", err)
	}

	return response.Metrics, nil
}

func (c *Client) FindMetricsPublic(ctx context.Context) (*model.DashboardMetricsPublic, error) {
	var response struct {
		Metrics *model.DashboardMetricsPublic `json:"metrics"`
	}

	if _, err := c.GET(ctx, "metrics/public", &response); err != nil {
		return nil, fmt.Errorf("find metrics: %w", err)
	}

	return response.Metrics, nil
}
