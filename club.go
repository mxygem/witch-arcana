package witcharcana

import "fmt"

// Clubs is a key value store of clubs with keys being the club's initials.
type Clubs struct {
	clubs map[string]*Club
	log   bool
}

// type Clubs map[string]*Club

// Club represents a particular club's data and the players that are currently members.
type Club struct {
	Name     string   `json:"name"`
	Location Location `json:"location,omitempty"`
	Players  Players  `json:"players,omitempty"`
}

func (cs *Clubs) All() map[string]*Club {
	return cs.clubs
}

func (cs *Clubs) Club(name string) (*Club, error) {
	if c := club(cs.clubs, name); c != nil {
		return c, nil
	}

	return nil, fmt.Errorf("no club %q found", name)
}

func club(cs map[string]*Club, name string) *Club {
	if c, ok := cs[name]; ok {
		return c
	}

	return nil
}

func maybeMakeClub(cs Clubs, name string) *Club {
	c, ok := cs.clubs[name]
	if ok {
		return c
	}

	nc := &Club{Name: name}
	cs.clubs[nc.Name] = nc

	return nc
}
