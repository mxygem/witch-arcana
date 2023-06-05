package witcharcana

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/pretty"
)

func Print(o any) error {
	b, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("marshalling data: %w", err)
	}

	opts := pretty.DefaultOptions
	opts.Width = 1000

	fmt.Println(string(pretty.PrettyOptions(b, opts)))

	return nil
}
