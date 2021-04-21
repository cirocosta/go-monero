package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/cirocosta/go-monero/pkg/levin"
)

var (
	fpath = flag.String("f", "resp.bin", "file location")
)

func run() error {
	flag.Parse()

	f, err := os.Open(*fpath)
	if err != nil {
		return fmt.Errorf("open %s: %w", *fpath, err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	ps, err := levin.NewPortableStorageFromBytes(b)
	if err != nil {
		return fmt.Errorf("portable storage from bytes: %w", err)
	}

	// if err := json.NewEncoder(os.Stdout).Encode(ps); err != nil {
	// 	return fmt.Errorf("encode: %w", err)
	// }

	for _, entry := range ps.Entries {
		if entry.Name != "local_peerlist_new" {
			continue
		}

		peerList, ok := entry.Value.(levin.Entries)
		if !ok {
			panic(fmt.Errorf("unexpected peerlist value type"))
		}

		for _, peer := range peerList {
			peerListAdr, ok := peer.Value.(levin.Entries)
			if !ok {
				panic(fmt.Errorf("unexpected value peerlist adr type"))
			}

			for _, adr := range peerListAdr {
				if adr.Name != "adr" {
					continue
				}

				addr, ok := adr.Value.(levin.Entries)
				if !ok {
					panic(fmt.Errorf("unexpected value type"))
				}

				for _, addrField := range addr {

					if addrField.Name != "addr" {
						continue
					}

					fields, ok := addrField.Value.(levin.Entries)
					if !ok {
						panic(fmt.Errorf("unexpected value type"))
					}

					var ip uint32
					var port uint16

					for _, field := range fields {
						if field.Name == "m_ip" {
							ip = field.Value.(uint32)
						}

						if field.Name == "m_port" {
							port = field.Value.(uint16)
						}
					}

					if ip != 0 && port != 0 {
						fmt.Printf("%s:%d\n", ipzify(ip), port)
					}
				}

			}
		}
	}

	return nil
}

func ipzify(ip uint32) string {
	result := make(net.IP, 4)
	result[0] = byte(ip)
	result[1] = byte(ip >> 8)
	result[2] = byte(ip >> 16)
	result[3] = byte(ip >> 24)

	return result.String()
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
