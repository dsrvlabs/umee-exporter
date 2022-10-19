package rpc

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	GetSlashWindowResponse struct {
		SlashWindow	string	`json:"window_progress"`
	}
)

func (c *RPCClient) GetSlashWindow(ctx context.Context) (*string, error) {
	body, err := c.rpcRequestGet(ctx, "/umee/oracle/v1/slash_window")
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed : %w", err)
	}

	var resp GetSlashWindowResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &resp.SlashWindow, nil
}
