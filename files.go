package witcharcana

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/tidwall/pretty"
)

func open(loc string) (*Clubs, error) {
	dat, err := os.ReadFile(loc)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}

	cs := &Clubs{}
	if len(dat) > 0 {
		if err = json.Unmarshal(dat, &cs.clubs); err != nil {
			return nil, fmt.Errorf("unmarshaling: %w", err)
		}
	}

	return cs, nil

	// dataFile, err := os.OpenFile(loc, os.O_CREATE|os.O_RDWR, 0644)
	// if err != nil {
	// 	return nil, fmt.Errorf("opening: %w", err)
	// }
	// defer dataFile.Close()

	// cs := clubs{}
	// if err := json.NewDecoder(dataFile).Decode(&cs); err != nil {
	// 	return nil, fmt.Errorf("decoding clubs: %w", err)
	// }

	// return cs, nil
}

// ReadCSV attempts to read in a csv file from the given location for bulk changes.
func ReadCSV(loc string) (Players, error) {
	inputFile, err := os.OpenFile(loc, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}
	defer inputFile.Close()

	ps := Players{}
	if err := gocsv.UnmarshalFile(inputFile, &ps); err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	return ps, nil
}

func (loc *Location) UnmarshalCSV(csv string) error {
	split := strings.Split(csv, ":")
	if len(split) == 0 {
		return nil
	}
	if len(split) != 2 {
		return fmt.Errorf("invalid location of %q received", csv)
	}

	x, err := strconv.Atoi(split[0])
	if err != nil {
		return fmt.Errorf("invalid integer %q received for location X value", split[0])
	}
	y, err := strconv.Atoi(split[1])
	if err != nil {
		return fmt.Errorf("invalid integer %q received for location Y value", split[1])
	}

	loc.X = x
	loc.Y = y

	return nil
}

func Save(loc string, cs *Clubs) error {
	b, err := json.Marshal(cs)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	p := pretty.Pretty(b)

	if err = os.WriteFile(loc, p, 0644); err != nil {
		return fmt.Errorf("writing: %w", err)
	}

	return nil
}
