package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var parser = flags.NewParser(&options, flags.Default)

type Options struct {
	Verbose bool   `short:"v" long:"verbose" description:"dump http requests and responses to stderr"`
	Address string `short:"a" long:"address" default:"http://xps.utxo.com.br:18089" description:"RPC server address"`
}

var options Options

func main() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
