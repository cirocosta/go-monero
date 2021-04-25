package main

import (
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var parser = flags.NewParser(&options, flags.Default)

type Options struct {
	Verbose        bool          `short:"v" env:"MONEROD_VERBOSE" long:"verbose" description:"dump http requests and responses to stderr"`
	Address        string        `short:"a" env:"MONEROD_ADDRESS" long:"address" description:"RPC server address" required:"true"`
	RequestTimeout time.Duration `short:"t" env:"MONEROD_TIMEOUT" long:"timeout" description:"request timeout" default:"10s"`
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
