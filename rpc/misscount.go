package rpc

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	GetMissCountResponse struct {
		MissCounter	string	`json:"miss_counter"`
	}
)

func (c *RPCClient) GetMissCount(ctx context.Context, valoperAddr string) (*string, error) {
	body, err := c.rpcRequestGet(ctx, "/umee/oracle/v1/validators/"+valoperAddr+"/miss")
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed : %w", err)
	}

	var resp GetMissCountResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &resp.MissCounter, nil
}
