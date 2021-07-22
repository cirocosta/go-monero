package wallet

import (
	"context"
	"fmt"
)

const (
	methodAutoRefresh   = "auto_refresh"
	methodCreateAddress = "create_address"
	methodGetAccounts   = "get_accounts"
	methodGetAddress    = "get_address"
	methodGetBalance    = "get_balance"
	methodGetHeight     = "get_height"
	methodRefresh       = "refresh"
)

func (c *Client) GetAccounts(
	ctx context.Context, params GetAccountsRequestParameters,
) (*GetAccountsResult, error) {
	resp := &GetAccountsResult{}

	if err := c.JSONRPC(ctx, methodGetAccounts, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetAddress(
	ctx context.Context, params GetAddressRequestParameters,
) (*GetAddressResult, error) {
	resp := &GetAddressResult{}

	if err := c.JSONRPC(ctx, methodGetAddress, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetBalance gets the balance of the wallet configured for the wallet rpc
// server.
//
func (c *Client) GetBalance(
	ctx context.Context, params GetBalanceRequestParameters,
) (*GetBalanceResult, error) {
	resp := &GetBalanceResult{}

	if err := c.JSONRPC(ctx, methodGetBalance, params, resp); err != nil {
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

func (c *Client) AutoRefresh(
	ctx context.Context, enable bool, period int64,
) (*AutoRefreshResult, error) {
	resp := &AutoRefreshResult{}

	params := map[string]interface{}{
		"enable": enable,
		"period": period,
	}
	if err := c.JSONRPC(ctx, methodAutoRefresh, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) Refresh(
	ctx context.Context, startHeight uint64,
) (*RefreshResult, error) {
	resp := &RefreshResult{}

	params := map[string]interface{}{
		"start_height": startHeight,
	}
	if err := c.JSONRPC(ctx, methodRefresh, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetHeight(ctx context.Context) (*GetHeightResult, error) {
	resp := &GetHeightResult{}

	if err := c.JSONRPC(ctx, methodGetHeight, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}
