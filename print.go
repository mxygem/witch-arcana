package witcharcana

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/pretty"
)

func Pretty(o any) ([]byte, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("marshalling data: %w", err)
	}

	opts := pretty.DefaultOptions
	opts.Width = 1000

	return pretty.PrettyOptions(b, opts), nil
}

func Print(o any) error {
	out, err := Pretty(o)
	if err != nil {
		return fmt.Errorf("formatting output for printing: %w", err)
	}

	fmt.Println(string(out))

	return nil
}
