package daemon

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type setBansCommand struct {
	Host     net.IP
	Duration time.Duration

	JSON bool
}

func (c *setBansCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-bans",
		Short: "ban another nodes",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")
	cmd.Flags().IPVar(&c.Host, "host",
		nil, "ip address (string format) of the host to ban")
	cmd.MarkFlagRequired("host")

	cmd.Flags().DurationVar(&c.Duration, "duration",
		24*time.Hour, "for how long this host should be banned for")

	return cmd
}

func (c *setBansCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	params := daemon.SetBansRequestParameters{
		Bans: []daemon.SetBansBan{
			{
				Host:    c.Host.String(),
				Ban:     true,
				Seconds: int64(c.Duration.Seconds()),
			},
		},
	}
	resp, err := client.SetBans(ctx, params)
	if err != nil {
		return fmt.Errorf("set bans: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *setBansCommand) pretty(v *daemon.SetBansResult) {
	fmt.Println(v.Status)
}

func init() {
	RootCommand.AddCommand((&setBansCommand{}).Cmd())
}
