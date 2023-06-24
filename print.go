package witcharcana

import (
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/tidwall/pretty"
)

func PrettyJSON(o any) ([]byte, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("marshalling data: %w", err)
	}

	opts := pretty.DefaultOptions
	opts.Width = 1000

	return pretty.PrettyOptions(b, opts), nil
}

func Print(o any) error {
	out, err := PrettyJSON(o)
	if err != nil {
		return fmt.Errorf("formatting output for printing: %w", err)
	}

	fmt.Println(string(out))

	return nil
}

func PrettyDiscord(w io.Writer, data *Clubs) error {
	tmpl, err := template.New("").ParseFiles("./templates/discord.tmpl")
	if err != nil {
		return fmt.Errorf("initiating template: %w", err)
	}

	if err := tmpl.ExecuteTemplate(w, "discord.tmpl", data.clubs); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}
