package wallet

import (
	"context"
	"fmt"
)

const (
	methodGetBalance    = "get_balance"
	methodCreateAddress = "create_address"
)

// GetBalance gets the balance of the wallet configured for the wallet rpc
// server.
//
func (c *Client) GetBalance(
	ctx context.Context,
) (*GetBalanceResult, error) {
	resp := &GetBalanceResult{}

	if err := c.JSONRPC(ctx, methodGetBalance, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) CreateAddress(
	ctx context.Context, accountIndex uint, count uint, label string,
) (*CreateAddressResult, error) {
	resp := &CreateAddressResult{}

	params := map[string]interface{}{
		"account_index": accountIndex,
		"label":         label,
		"count":         count,
	}
	if err := c.JSONRPC(ctx, methodCreateAddress, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}
