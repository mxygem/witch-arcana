package witcharcana

import (
	"fmt"
	"log"
)

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
	c := &Club{
		Name: name,
	}

	if x > 0 && y > 0 {
		c.Location = &Location{X: x, Y: y}
	}

	return c
}

func (cs *Clubs) LoadData(filename string) error {
	csd, err := loadData(filename)
	if err != nil {
		return fmt.Errorf("loading data: %w", err)
	}

	if cs.log {
		log.Printf("found %d clubs\n", len(csd.clubs))
		for k, v := range csd.clubs {
			log.Printf("club %q with %d players\n", k, len(v.Players))
			for i, p := range v.Players {
				log.Printf("\tp %d:%+v\n", i, p)
			}
		}
	}

	cs.clubs = csd.clubs

	return nil
}

func loadData(filename string) (*Clubs, error) {
	csd, err := open(filename)
	if err != nil {
		return nil, fmt.Errorf("loading club data from file: %q: %w", filename, err)
	}

	return csd, nil
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

func (cs *Clubs) CreateClub(c *Club) error {
	// todo: return newly created club
	return createClub(cs, c)
}

func createClub(cs *Clubs, c *Club) error {
	if c.Name == "" {
		return fmt.Errorf("club name required")
	}

	if oc := club(cs.clubs, c.Name); oc != nil {
		return fmt.Errorf("club %q already exists", c.Name)
	}

	cs.clubs[c.Name] = c

	return nil
}

func (cs *Clubs) UpdateClub(uc *Club) (*Club, error) {
	c, err := updateClub(cs, uc)
	if err != nil {
		return nil, fmt.Errorf("updating club: %w", err)
	}

	return c, nil
}

func updateClub(cs *Clubs, uc *Club) (*Club, error) {
	if uc.Name == "" {
		return nil, fmt.Errorf("updated club information must contain a name")
	}
	if uc.Location.X == 0 || uc.Location.Y == 0 {
		return nil, fmt.Errorf("no location values can be zero. got x: %d y: %d", uc.Location.X, uc.Location.Y)
	}

	c := club(cs.clubs, uc.Name)
	if c == nil {
		return nil, fmt.Errorf("no club %q found", uc.Name)
	}

	if c.Location == nil {
		l := &Location{X: uc.Location.X, Y: uc.Location.Y}
		c.Location = l
	} else {
		if c.Location.X != uc.Location.X {
			c.Location.X = uc.Location.X
		}
		if c.Location.Y != uc.Location.Y {
			c.Location.Y = uc.Location.Y
		}
	}

	return c, nil
}

func (cs *Clubs) RemoveClub(name string) error {
	if err := removeClub(cs, name); err != nil {
		return fmt.Errorf("removing club: %w", err)
	}

	return nil
}

func removeClub(cs *Clubs, name string) error {
	if name == "" {
		return fmt.Errorf("club name required")
	}

	c := club(cs.clubs, name)
	if c == nil {
		return fmt.Errorf("no club %q found", name)
	}

	delete(cs.clubs, name)

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
