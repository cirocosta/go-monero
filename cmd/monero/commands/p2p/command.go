package p2p

import (
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "p2p",
	Short: "execute p2p commands against a monero node",
}
