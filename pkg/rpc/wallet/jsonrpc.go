package wallet

import (
	"context"
	"fmt"
)

const (
	methodGetBalance = "get_balance"
)

// GetBalance gets the balance of the wallet configured for the wallet rpc
// server.
//
func (c *Client) GetBalance(ctx context.Context) (*GetBalanceResult, error) {
	resp := &GetBalanceResult{}

	if err := c.JSONRPC(ctx, methodGetBalance, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}
