package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
	"golang.org/x/net/proxy"

	"github.com/cirocosta/go-monero/pkg/crawler"
	"github.com/cirocosta/go-monero/pkg/levin"
)

type CrawlCommand struct {
	Ip      string        `long:"ip"      default:"198.98.116.72" description:"p2p address"`
	Port    uint16        `long:"port"    default:"18080"         description:"p2p port"`
	Timeout time.Duration `long:"timeout" default:"20s"           description:"maximum execution time"`
	Output  string        `long:"output"  default:"nodes.csv"     description:"file to write peers to"`

	Proxy         string `long:"proxy" short:"x" description:"socks5 proxy addr"`
	GeoIpDatabase string `long:"geo-ip-db" description:"fpath of a mmdb geoip file"`
}

func (c *CrawlCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	f, err := os.Create(c.Output)
	if err != nil {
		return fmt.Errorf("create '%s': %w", c.Output, err)
	}
	defer f.Close()

	opts := []levin.ClientOption{}

	if c.Proxy != "" {
		dialer, err := proxy.SOCKS5("tcp", c.Proxy, nil, nil)
		if err != nil {
			return fmt.Errorf("socks5 '%s': %w", c.Proxy, err)
		}

		contextDialer, ok := dialer.(proxy.ContextDialer)
		if !ok {
			panic("can't cast proxy dialer to proxy context dialer")
		}

		opts = append(opts, levin.WithContextDialer(contextDialer))
	}

	ccrawler := crawler.NewCrawler(opts...)
	ccrawler.TryPutForVisit(&levin.Peer{
		Ip: c.Ip, Port: c.Port,
	})

	processingFuncs := []func(node *crawler.VisitedPeer) string{
		func(node *crawler.VisitedPeer) string {
			line := node.Addr() + ","
			if node.Error != nil {
				line += node.Error.Error()
			}

			return line
		},
	}

	if c.GeoIpDatabase != "" {
		db, err := geoip2.Open(c.GeoIpDatabase)
		if err != nil {
			return fmt.Errorf("geoip open: %w", err)
		}
		defer db.Close()

		processingFuncs = append(processingFuncs, func(node *crawler.VisitedPeer) string {
			ip := net.ParseIP(node.Ip())

			record, err := db.Country(ip)
			if err != nil {
				panic(fmt.Errorf("db city '%s': %w", ip, err))
			}

			return record.Country.Names["en"]
		})
	}

	go func() {
		for node := range ccrawler.C {
			columns := []string{}

			for _, f := range processingFuncs {
				columns = append(columns, f(node))
			}

			if _, err := f.WriteString(strings.Join(columns, ",") + "\n"); err != nil {
				panic(err)
			}
		}
	}()

	if _, err := ccrawler.Run(ctx); err != nil {
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
