package rpc

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	GetPrevoteSubmitResponse struct {
		Prevote struct {
			Hash		string	`json:"hash"`
			Voter		string	`json:"voter"`
			SubmitBlock	string	`json:"submit_block"`
		} `json:"aggregate_prevote"`
	}
)

func (c *RPCClient) GetPrevoteSubmit(ctx context.Context, valoperAddr string) (*string, error) {
	body, err := c.rpcRequestGet(ctx, "/umee/oracle/v1/validators/"+valoperAddr+"/aggregate_prevote")
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed : %w", err)
	}

	var resp GetPrevoteSubmitResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &resp.Prevote.SubmitBlock, nil
}
