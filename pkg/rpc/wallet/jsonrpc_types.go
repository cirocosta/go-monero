package wallet

type GetAccountsRequestParameters struct {
	Tag            string `json:"tag,omitempty"`
	StrictBalances bool   `json:"strict_balances,omitempty"`
}

type GetAccountsResult struct {
	SubaddressAccounts []struct {
		AccountIndex    uint   `json:"account_index"`
		Balance         uint64 `json:"balance"`
		BaseAddress     string `json:"base_address"`
		Label           string `json:"label"`
		Tag             string `json:"tag"`
		UnlockedBalance uint64 `json:"unlocked_balance"`
	} `json:"subaddress_accounts"`

	TotalBalance         uint64 `json:"total_balance"`
	TotalUnlockedBalance uint64 `json:"total_unlocked_balance"`
}

type GetAddressRequestParameters struct {
	AccountIndex   uint   `json:"account_index"`
	AddressIndices []uint `json:"address_indices"`
}

type GetAddressResult struct {
	Address   string `json:"address"`
	Addresses []struct {
		Address      string `json:"address"`
		AddressIndex uint   `json:"address_index"`
		Label        string `json:"label"`
		Used         bool   `json:"used"`
	} `json:"addresses"`
}

type GetBalanceRequestParameters struct {
	AccountIndex   uint   `json:"account_index"`
	AddressIndices []uint `json:"address_indices"`
	AllAccounts    bool   `json:"all_accounts"`
	Strict         bool   `json:"strict"`
}

type GetBalanceResult struct {
	// Balance is the total balance of the current monero-wallet-rpc in
	// session.
	//
	Balance uint64 `json:"balance"`

	// BlocksToUnlock indicates how many blocks are necessary for all the
	// funds to be unclocked.
	//
	BlocksToUnlock uint `json:"blocks_to_unlock"`

	// MultisigImportNeeded is True if importing multisig data is needed
	// for returning a correct balance
	//
	MultisigImportNeeded bool `json:"multisig_import_needed"`

	// PerSubaddress is an array of subaddress information; Balance
	// information for each subaddress in an account.
	//
	PerSubaddress []SubAddress `json:"per_subaddress"`

	// TimeToUnlock TODO
	//
	TimeToUnlock int `json:"time_to_unlock"`

	// UnlockedBalance TODO
	//
	UnlockedBalance int64 `json:"unlocked_balance"`
}

type SubAddress struct {
	// AccountIndex is the index of the account.
	//
	AccountIndex uint `json:"account_index"`

	// Address at this index. Base58 representation of the public
	// keys.
	//
	Address string `json:"address"`

	// AddressIndex  is the index of the subaddress in the account.
	//
	AddressIndex uint `json:"address_index"`

	// Balance is the balance for the subaddress.
	//
	Balance uint64 `json:"balance"`

	// BlocksToUnlock TODO
	//
	BlocksToUnlock uint `json:"blocks_to_unlock"`

	// Label TODO
	//
	Label string `json:"label"`

	// NumUnspentOutputs TODO
	//
	NumUnspentOutputs uint `json:"num_unspent_outputs"`

	// TimeToUnlock TODO
	//
	TimeToUnlock uint `json:"time_to_unlock"`

	// UnlockedBalance TODO
	//
	UnlockedBalance int64 `json:"unlocked_balance"`
}

type CreateAddressResult struct {
	Address        string   `json:"address"`
	AddressIndex   uint     `json:"address_index"`
	AddressIndices []uint   `json:"address_indices"`
	Addresses      []string `json:"addresses"`
}

type RefreshResult struct {
	BlocksFetched uint64 `json:"blocks_fetched"`
	ReceivedMoney bool   `json:"received_money"`
}

type AutoRefreshResult struct {
}

type GetHeightResult struct {
	Height uint64 `json:"height"`
}
