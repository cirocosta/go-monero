package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gosuri/uitable"
)

// displayResponseInJSONFormat pushes to stdout a pretty printed representation
// of a given value `v`.
//
func displayResponseInJSONFormat(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")

	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

func newTable() *uitable.Table {
	table := uitable.New()

	table.MaxColWidth = 80
	table.Wrap = true

	return table
}
