package levin_test

import (
	"fmt"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/cirocosta/go-monero/pkg/levin"
)

func TestLevin(t *testing.T) {

	spec.Run(t, "Request", func(t *testing.T, when spec.G, it spec.S) {

		it("new", func() {
			req := levin.NewRequest(levin.CommandPing, []byte("test"))
			fmt.Println(req.Bytes())

		})

	}, spec.Report(report.Terminal{}), spec.Parallel(), spec.Random())
}
