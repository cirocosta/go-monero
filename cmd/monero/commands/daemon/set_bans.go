package daemon

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type setBansCommand struct {
	Host     net.IP
	Duration time.Duration
	Filepath string

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

	cmd.Flags().StringVarP(&c.Filepath, "filepath", "f",
		"", "location of a csv file containing <host>,<period> "+
			"entries to ban")

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
		Bans: []daemon.SetBansBan{},
	}

	if c.Host != nil {
		params.Bans = append(params.Bans, daemon.SetBansBan{
			Host:    c.Host.String(),
			Ban:     true,
			Seconds: int64(c.Duration.Seconds()),
		})
	}

	if c.Filepath != "" {
		bansFromFile, err := c.bansFromFilepath()
		if err != nil {
			return fmt.Errorf("bans from file: %w", err)
		}

		params.Bans = append(params.Bans, bansFromFile...)
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

func (c *setBansCommand) bansFromFilepath() ([]daemon.SetBansBan, error) {
	f, err := os.Open(c.Filepath)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	entries, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv read all: %w", err)
	}

	bans := []daemon.SetBansBan{}
	for _, entry := range entries {
		if len(entry) != 2 {
			return nil, fmt.Errorf(
				"expected 2 fields in entry, got %d",
				len(entry),
			)
		}

		host, durationStr := entry[0], entry[1]
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return nil, fmt.Errorf("parse duration: %w", err)
		}

		bans = append(bans, daemon.SetBansBan{
			Host:    host,
			Seconds: int64(duration.Seconds()),
			Ban:     true,
		})
	}

	return bans, nil
}

// nolint:forbidigo
func (c *setBansCommand) pretty(v *daemon.SetBansResult) {
	fmt.Println(v.Status)
}

func init() {
	RootCommand.AddCommand((&setBansCommand{}).Cmd())
}
