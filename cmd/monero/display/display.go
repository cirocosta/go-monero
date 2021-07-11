package display

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gosuri/uitable"
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

	table.MaxColWidth = 80
	table.Wrap = true

	return table
}
