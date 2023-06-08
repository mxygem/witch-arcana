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
	Name     string    `json:"name"`
	Location *Location `json:"location,omitempty"`
	Players  Players   `json:"players,omitempty"`
}

func NewClubs(log bool) *Clubs {
	return &Clubs{log: log}
}

func NewClub(name string, x, y int) *Club {
	return &Club{
		Name:     name,
		Location: &Location{X: x, Y: y},
	}
}

func (cs *Clubs) LoadData(filename string) error {
	csd, err := open(filename)
	if err != nil {
		return fmt.Errorf("loading club data from file: %q: %w", filename, err)
	}

	if cs.log {
		fmt.Printf("found %d clubs\n", len(csd.clubs))
		for k, v := range csd.clubs {
			fmt.Printf("club %q with %d players\n", k, len(v.Players))
			for i, p := range v.Players {
				fmt.Printf("\tp %d:%+v\n", i, p)
			}
		}
	}

	cs.clubs = csd.clubs
	return nil
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

func (cs *Clubs) CreateClub(c *Club) error {
	return nil
}

func createClub(cs *Clubs, c *Club) error {
	return nil
}

func (cs *Clubs) UpdateClub(uc *Club) error {
	return nil
}

func updateClub(oc *Club, uc *Club) *Club {
	return nil
}

func club(cs map[string]*Club, name string) *Club {
	if c, ok := cs[name]; ok {
		return c
	}

	return nil
}

func (cs *Clubs) RemoveClub(name string) error {
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
