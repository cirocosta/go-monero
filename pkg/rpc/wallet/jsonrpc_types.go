package wallet

type GetBalanceResult struct {
	// Balance is the total balance of the current monero-wallet-rpc in session.
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
	PerSubaddress []struct {
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
	} `json:"per_subaddress"`

	// TimeToUnlock TODO
	//
	TimeToUnlock int `json:"time_to_unlock"`

	// UnlockedBalance TODO
	//
	UnlockedBalance int64 `json:"unlocked_balance"`
}
