package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/commands/daemon"
	"github.com/cirocosta/go-monero/cmd/monero/commands/p2p"
	"github.com/cirocosta/go-monero/cmd/monero/commands/wallet"
)

var (
	version = "dev"
	commit  = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "monero",
	Short: "Daemon, Wallet, and p2p command line monero CLI",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of this cli",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(version, commit)
	},
}

func init() {
	rootCmd.AddCommand(daemon.RootCommand)
	rootCmd.AddCommand(wallet.RootCommand)
	rootCmd.AddCommand(p2p.RootCommand)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
