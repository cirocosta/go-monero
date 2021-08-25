package address

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/pkg/monero"
)

type generateCommand struct {
	passphrase  string
	networkName string
}

func (c *generateCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "create a new single-key seed",
		RunE:  c.RunE,
	}

	cmd.Flags().StringVar(&c.passphrase, "passphrase", "",
		"string to encode and use as offset for the seed")

	cmd.Flags().StringVar(&c.networkName, "network", "mainnet",
		"network that the addresses should be used for "+
			c.networkOptions())

	return cmd
}

func (c *generateCommand) RunE(_ *cobra.Command, _ []string) error {
	privateKey, err := c.privateKey()
	if err != nil {
		return fmt.Errorf("private key: %w", err)
	}

	network, err := c.network()
	if err != nil {
		return fmt.Errorf("network: %w", err)
	}

	c.pretty(monero.NewSeed(privateKey, monero.WithNetwork(network)))
	return nil
}

func (c *generateCommand) pretty(seed *monero.Seed) {
	c.prettyMnemonic(seed)
	c.prettyKeys(seed)
}

func (c *generateCommand) prettyKeys(seed *monero.Seed) {
	table := display.NewTable()
	defer fmt.Println(table)

	table.AddRow("Primary Address:", seed.PrimaryAddress())
	table.AddRow("Private Spend Key:",
		hex.EncodeToString(seed.PrivateSpendKey()))
	table.AddRow("Private View Key:",
		hex.EncodeToString(seed.PrivateViewKey()))
	table.AddRow("Public Spend Key:",
		hex.EncodeToString(seed.PublicSpendKey()))
	table.AddRow("Public View Key:",
		hex.EncodeToString(seed.PublicViewKey()))
}

func (c *generateCommand) prettyMnemonic(seed *monero.Seed) {
	table := display.NewTable()
	defer fmt.Println(table)

	table.Separator = "  "

	mnemonic := seed.Mnemonic()

	table.AddRow(c.row("Mnemonic:", mnemonic[0:4]...)...)
	table.AddRow(c.row("", mnemonic[4:8]...)...)
	table.AddRow(c.row("", mnemonic[8:12]...)...)
	table.AddRow(c.row("", mnemonic[12:16]...)...)
	table.AddRow(c.row("", mnemonic[16:20]...)...)
	table.AddRow(c.row("", mnemonic[20:24]...)...)
	table.AddRow(c.row("", mnemonic[24])...)
	table.AddRow("")
}

func (c *generateCommand) row(key string, values ...string) []interface{} {
	res := []interface{}{key}
	for _, v := range values {
		res = append(res, v)
	}

	return res
}

func (c *generateCommand) privateKey() ([]byte, error) {
	v := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, v)
	if err != nil {
		return nil, fmt.Errorf("read full: %w", err)
	}

	return v, nil
}

func (c *generateCommand) networkOptions() string {
	strs := []string{}

	for _, network := range []monero.Network{
		monero.NetworkMainnet,
		monero.NetworkTestnet,
		monero.NetworkStagenet,
		monero.NetworkFakechain,
	} {
		strs = append(strs, string(network))
	}

	return "(" + strings.Join(strs, ",") + ")"
}

func (c *generateCommand) network() (monero.Network, error) {
	switch c.networkName {
	case string(monero.NetworkMainnet):
		return monero.NetworkMainnet, nil
	case string(monero.NetworkTestnet):
		return monero.NetworkTestnet, nil
	case string(monero.NetworkStagenet):
		return monero.NetworkStagenet, nil
	case string(monero.NetworkFakechain):
		return monero.NetworkFakechain, nil
	}

	err := fmt.Errorf("unknown network %s", c.network)
	return monero.NetworkFakechain, err
}

func init() {
	RootCommand.AddCommand((&generateCommand{}).Cmd())
}
