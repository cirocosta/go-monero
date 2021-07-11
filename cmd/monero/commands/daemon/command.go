package daemon

import (
	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/options"
)

var RootCommand = &cobra.Command{
	Use:   "daemon",
	Short: "execute remote procedure calls against a monero node",
}

func init() {
	options.Bind(RootCommand)
}
