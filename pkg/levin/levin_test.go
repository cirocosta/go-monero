package levin_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/assert"

	"github.com/cirocosta/go-monero/pkg/levin"
)

func TestLevin(t *testing.T) {

	spec.Run(t, "newheader", func(t *testing.T, when spec.G, it spec.S) {

		it("request", func() {
			bytes := levin.NewHeader(levin.CommandPing, 1).Bytes()

			assert.ElementsMatch(t, bytes, []byte{
				0x01, 0x21, 0x01, 0x01, // signature
				0x01, 0x01, 0x01, 0x01,

				0x01, 0x00, 0x00, 0x00, // length       -- 0 for a ping msg
				0x00, 0x00, 0x00, 0x00,

				0x01, // expects response               -- `true` bool

				0x03, 0x10, 0x00, 0x00, // command	-- 1003 for ping

				0x00, 0x00, 0x00, 0x00, // return code	-- 0 for requests

				0x01, 0x00, 0x00, 0x00, // flags	-- Q(1st lsb) set for req

				0x00, 0x00, 0x00, 0x00, // end
			})
		})

	}, spec.Report(report.Terminal{}), spec.Parallel(), spec.Random())
}
