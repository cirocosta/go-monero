package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cirocosta/go-monero/pkg/crawler"
	"github.com/cirocosta/go-monero/pkg/levin"
)

type CrawlCommand struct {
	Ip      string        `long:"ip"      default:"xps.utxo.com.br" description:"p2p address"`
	Port    uint16        `long:"port"    default:"18080"           description:"p2p port"`
	Timeout time.Duration `long:"timeout" default:"20s"             description:"maximum execution time"`
	Output  string        `long:"output"  default:"nodes.csv"       description:"file to write peers to"`
}

func (c *CrawlCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	f, err := os.Create(c.Output)
	if err != nil {
		return fmt.Errorf("create '%s': %w", c.Output, err)
	}
	defer f.Close()

	crawler := crawler.NewCrawler()
	crawler.TryPutForVisit(&levin.Peer{
		Ip: c.Ip, Port: c.Port,
	})

	go func() {
		for node := range crawler.C {
			line := node.Addr() + ","
			if node.Error != nil {
				line += node.Error.Error()
			}

			line += "\n"

			if _, err := f.WriteString(line); err != nil {
				panic(err)
			}
		}
	}()

	if _, err := crawler.Run(ctx); err != nil {
		return fmt.Errorf("crawler run: %w", err)
	}

	return nil
}

func init() {
	parser.AddCommand("crawl",
		"Crawl over the network to find all peers",
		"Crawl over the network to find all peers",
		&CrawlCommand{},
	)
}
