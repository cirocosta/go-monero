package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

// configure global logging
//
func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var options RPCOptions
var parser = flags.NewParser(&options, flags.Default)

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
