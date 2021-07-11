package daemon

import (
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
	"github.com/gosuri/uitable"
)

func prettyBlockHeader(table *uitable.Table, header daemon.BlockHeader) {
	table.AddRow("Block Size:", header.BlockSize)
	table.AddRow("Block Weight:", header.BlockWeight)
	table.AddRow("Cumulative Difficulty:", header.CumulativeDifficulty)
	table.AddRow("Cumulative Difficulty Top64:", header.CumulativeDifficultyTop64)
	table.AddRow("Depth:", header.Depth)
	table.AddRow("Difficulty:", header.Difficulty)
	table.AddRow("Difficulty Top64:", header.DifficultyTop64)
	table.AddRow("Hash:", header.Hash)
	table.AddRow("Height:", header.Height)
	table.AddRow("Long Term Weight:", header.LongTermWeight)
	table.AddRow("Major Version:", header.MajorVersion)
	table.AddRow("Miner Transaction Hash:", header.MinerTxHash)
	table.AddRow("Minor Version:", header.MinorVersion)
	table.AddRow("Nonce:", header.Nonce)
	table.AddRow("Number of Transactions:", header.NumTxes)
	table.AddRow("Orphan Status:", header.OrphanStatus)
	table.AddRow("Proof-of-Work Hash:", header.PowHash)
	table.AddRow("Previous Hash:", header.PrevHash)
	table.AddRow("Reward:", header.Reward)
	table.AddRow("Timestamp:", header.Timestamp)
	table.AddRow("WIDE Cumulative Difficulty:", header.WideCumulativeDifficulty)
	table.AddRow("Wide Difficulty:", header.WideDifficulty)
}
