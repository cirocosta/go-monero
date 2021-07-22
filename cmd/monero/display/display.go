package display

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gosuri/uitable"

	"github.com/cirocosta/go-monero/pkg/constant"
)

// JSON pushes to stdout a pretty printed representation of a given value `v`.
//
func JSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")

	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

// NewTable instantiates a new table instance that already has pre-defined
// options set so it's consistent across all pretty prints of the commands.
//
func NewTable() *uitable.Table {
	table := uitable.New()

	table.MaxColWidth = 160

	return table
}

func XMR(v uint64) string {
	return fmt.Sprintf("%.2f ɱ", float64(v)/float64(constant.XMR))
}

func ShortenAddress(addr string) string {
	if len(addr) < 10 {
		return addr
	}

	return addr[:5] + ".." + addr[len(addr)-5:]
}
