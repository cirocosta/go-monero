package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cirocosta/go-monero/pkg/crawler"
	"github.com/cirocosta/go-monero/pkg/levin"
)

type CrawlCommand struct {
	Ip      string        `long:"ip"      default:"xps.utxo.com.br" description:"p2p address"`
	Port    uint16        `long:"port"    default:"18080"           description:"p2p port"`
	Timeout time.Duration `long:"timeout" default:"20s"             description:"maximum execution time"`
}

func (c *CrawlCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	crawler := crawler.NewCrawler()

	crawler.TryPutForVisit(&levin.Peer{
		Ip: c.Ip, Port: c.Port,
	})

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
