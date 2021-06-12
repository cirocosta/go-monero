package main

import (
	"fmt"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

type GetBansCommand struct {
	Json bool `long:"json" description:"output the raw json response from the endpoint"`
}

func (c *GetBansCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBans(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.Json {
		return displayResponseInJSONFormat(resp)
	}

	c.display(resp)
	return nil
}

func (c *GetBansCommand) display(v *rpc.GetBansResult) {
	table := newTable()

	for _, ban := range v.Bans {
		table.AddRow("Host:", ban.Host)
		table.AddRow("IP:", ban.IP)
		table.AddRow("Seconds:", ban.Seconds)
		table.AddRow("")
	}

	fmt.Println(table)
}

func init() {
	parser.AddCommand("get-bans",
		"Get bans",
		"Get bans",
		&GetBansCommand{},
	)
}
