package witcharcana

import (
	"fmt"
	"log"
)

// Players is a collection of players.
type Players []*Player

// Player represents a player and various data about them.
type Player struct {
	Name     string    `json:"name" csv:"name"`
	Location *Location `json:"location,omitempty" csv:"location"`
	InHive   bool      `json:"in_hive,omitempty" csv:"in_hive"`
	Level    int       `json:"level,omitempty" csv:"level,lvl"`
	Might    int64     `json:"might,omitempty" csv:"might"`
	Club     string    `json:"club,omitempty" csv:"club"`
}

func NewPlayer(name, clubName string, level, x, y int) *Player {
	p := &Player{
		Name:  name,
		Club:  clubName,
		Level: level,
	}

	if x > 0 && y > 0 {
		p.Location = &Location{X: x, Y: y}
	}

	return p
}

// Player returns a given player if found.
func (cs *Clubs) Player(name string) (*Player, error) {
	if _, p := player(cs, name); p != nil {
		return p, nil
	}

	return nil, fmt.Errorf("player %q not found", name)
}

func player(cs *Clubs, name string) (int, *Player) {

	var c Club
	var p Player
	for _, club := range cs.clubs {
		for _, pl := range club.Players {
			if pl.Name != name {
				continue
			}

			p = *pl
			p.Club = club.Name
			c = *club

			break
		}
	}

	if &c == nil || &p == nil {
		return -1, nil
	}

	if cs.db != nil {
		fp, err := cs.db.Player(c.Name, p.Name)
		if err != nil {
			log.Printf("getting player from db: %v", err)
			return -1, nil
		}

		return 0, fp
	}

	return -1, nil
}

// CreatePlayer creates a new player.
func (cs *Clubs) CreatePlayer(clubName string, np *Player) (*Player, error) {
	fmt.Printf("cs.CreatePlayer: clubName %q, np %q\n", clubName, np.Name)
	c, err := cs.Club(clubName)
	if err != nil {
		return nil, fmt.Errorf("getting club: %w", err)
	}
	if c == nil {
		return nil, fmt.Errorf("club %q not found", clubName)
	}

	if err := createPlayer(cs, c, np); err != nil {
		return nil, fmt.Errorf("creating player: %w", err)
	}

	_, p := player(cs, np.Name)
	fmt.Printf("found player after create: %+v\n", p)

	return p, nil
}

func createPlayer(cs *Clubs, c *Club, np *Player) error {
	fmt.Printf("checking if player p.Name exists: %q\n", np.Name)
	for _, p := range c.Players {
		if p.Name == np.Name {
			return fmt.Errorf("player %q already exists in club %q", np.Name, c.Name)
		}
	}

	if np.Club != "" {
		np.Club = ""
	}

	c.Players = append(c.Players, np)

	if cs.db != nil {
		if err := cs.db.Update(c); err != nil {
			return fmt.Errorf("updating db: %w", err)
		}
	}

	return nil
}

// RemovePlayer completely removes a player.
func (cs *Clubs) RemovePlayer(playerName string) error {
	if err := removePlayer(cs, playerName); err != nil {
		return fmt.Errorf("removing player: %w", err)
	}

	return nil
}

func removePlayer(clubs *Clubs, playerName string) error {
	pos, player := player(clubs, playerName)
	if player == nil {
		return fmt.Errorf("player %q does not exist", playerName)
	}

	c := club(clubs.clubs, player.Club)
	c.Players = append(c.Players[:pos], c.Players[pos+1:]...)

	return nil
}

// MovePlayer moves a player from one club to another.
func (cs *Clubs) MovePlayer(playerName, newClubName string) (*Player, error) {
	player, err := movePlayer(cs, newClubName, playerName)
	if err != nil {
		return nil, fmt.Errorf("unable to move player %q to %q: %w", playerName, newClubName, err)
	}

	return player, nil
}

func movePlayer(cs *Clubs, newClubName, playerName string) (*Player, error) {
	_, p := player(cs, playerName)
	if p == nil {
		return nil, fmt.Errorf("player %q does not exist", playerName)
	}

	c, err := cs.Club(newClubName)
	if err != nil {
		return nil, fmt.Errorf("getting new club: %w", err)
	}

	if err := createPlayer(cs, c, p); err != nil {
		return nil, fmt.Errorf("creating player in new club: %w", err)
	}

	if err := removePlayer(cs, p.Name); err != nil {
		return nil, fmt.Errorf("removing player from old club: %w", err)
	}

	_, p = player(cs, p.Name)
	if p == nil {
		return nil, fmt.Errorf("could not find player %q after move", p.Name)
	}

	return p, nil
}

func (cs *Clubs) UpdatePlayer(p *Player) (*Player, error) {
	n, fp := player(cs, p.Name)
	if n < 0 {
		return nil, fmt.Errorf("player %q not found", p.Name)
	}

	up := updatePlayer(fp, p)

	c := club(cs.clubs, p.Club)
	c.Players[n] = up

	return up, nil
}

func updatePlayer(p *Player, up *Player) *Player {
	p.Club = ""

	if p.Location == nil && (up.Location.X > 0 && up.Location.Y > 0) {
		l := &Location{X: up.Location.X, Y: up.Location.Y}
		p.Location = l
	} else {
		if p.Location.X != up.Location.X {
			p.Location.X = up.Location.X
		}
		if p.Location.Y != up.Location.Y {
			p.Location.Y = up.Location.Y
		}
	}
	if up.Level != 0 && up.Level != p.Level {
		p.Level = up.Level
	}
	if up.Location.X != 0 && up.Location.X != p.Location.X {
		p.Location.X = up.Location.X
	}
	if up.Location.Y != 0 && up.Location.Y != p.Location.Y {
		p.Location.Y = up.Location.Y
	}
	if up.InHive != p.InHive {
		p.InHive = up.InHive
	}

	return p
}

// BulkUpdatePlayers creates/updates players based on data read in from a csv.
func (cs *Clubs) BulkUpdatePlayers(ps Players) error {
	return bulkUpdatePlayers(cs, ps)
}

func bulkUpdatePlayers(cs *Clubs, ps Players) error {
	// todo: add in moving players
	for _, np := range ps {
		c := maybeMakeClub(cs, np.Club)
		n, p := player(cs, np.Name)
		if n < 0 {
			if err := createPlayer(cs, c, np); err != nil {
				return fmt.Errorf("bulk update: creating player: %w", err)
			}
			continue
		}

		up := updatePlayer(p, np)
		if up == nil {
			return fmt.Errorf("bulk update: failed to update player: %q", p.Name)
		}
		c.Players[n] = up
	}

	return nil
}
