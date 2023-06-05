package witcharcana

import "fmt"

func (cs Clubs) Club(name string) (*Club, error) {
	if c := club(cs, name); c != nil {
		return c, nil
	}

	return nil, fmt.Errorf("no club %q found", name)
}

func club(cs Clubs, name string) *Club {
	if c, ok := cs[name]; ok {
		return c
	}

	return nil
}
