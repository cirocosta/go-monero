package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/gosuri/uitable"
	"github.com/oschwald/geoip2-golang"
)

var (
	file = flag.String("file", "", "mmdb geoip file")
	addr = flag.String("addr", "", "address to figure location out")
)

func run() error {
	if *file == "" {
		return fmt.Errorf("mmdb file must be specified")
	}

	if *addr == "" {
		return fmt.Errorf("non-empty address must be specified")
	}

	db, err := geoip2.Open(*file)
	if err != nil {
		return fmt.Errorf("geoip open: %w", err)
	}
	defer db.Close()

	ip := net.ParseIP(*addr)
	record, err := db.City(ip)
	if err != nil {
		return fmt.Errorf("db city: %w", err)
	}

	table := uitable.New()
	table.AddRow("Continent:", record.Continent.Names["en"])
	table.AddRow("Country:", record.Country.Names["en"])
	table.AddRow("City:", record.City.Names["en"])
	table.AddRow("Coordinates:", fmt.Sprintf("(%f,%f)", record.Location.Longitude, record.Location.Latitude))

	fmt.Println(table)

	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}
