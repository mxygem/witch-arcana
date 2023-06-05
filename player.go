package witcharcana

import (
	"fmt"
)

// Player returns a given player if found.
func (cs Clubs) Player(name string) (*Player, error) {
	if _, p := player(cs, name); p != nil {
		return p, nil
	}

	return nil, fmt.Errorf("player %q not found", name)
}

func player(cs Clubs, name string) (int, *Player) {
	for _, club := range cs {
		for i, p := range club.Players {
			if p.Name != name {
				continue
			}

			op := *p
			op.Club = club.Name

			return i, &op
		}
	}

	return -1, nil
}

// CreatePlayer creates a new player.
func (cs Clubs) CreatePlayer(clubName string, newPlayer *Player) (*Player, error) {
	c := club(cs, clubName)
	if c == nil {
		return nil, fmt.Errorf("club %q not found", clubName)
	}

	if err := createPlayer(cs[clubName], newPlayer); err != nil {
		return nil, fmt.Errorf("creating player: %w", err)
	}

	_, p := player(cs, newPlayer.Name)

	return p, nil
}

func createPlayer(club *Club, player *Player) error {
	for _, p := range club.Players {
		if p.Name == player.Name {
			return fmt.Errorf("player %q already exists in club %q", player.Name, club.Name)
		}
	}

	if player.Club != "" {
		player.Club = ""
	}
	club.Players = append(club.Players, player)

	return nil
}

// RemovePlayer completely removes a player.
func (cs Clubs) RemovePlayer(playerName string) error {
	if err := removePlayer(cs, playerName); err != nil {
		return fmt.Errorf("removing player: %w", err)
	}

	return nil
}

func removePlayer(clubs Clubs, playerName string) error {
	pos, player := player(clubs, playerName)
	if player == nil {
		return fmt.Errorf("player %q does not exist", playerName)
	}

	c := club(clubs, player.Club)
	c.Players = append(c.Players[:pos], c.Players[pos+1:]...)

	return nil
}

// MovePlayer moves a player from one club to another.
func (cs Clubs) MovePlayer(playerName, newClubName string) (*Player, error) {
	player, err := movePlayer(cs, newClubName, playerName)
	if err != nil {
		return nil, fmt.Errorf("unable to move player %q to %q: %w", playerName, newClubName, err)
	}

	return player, nil
}

func movePlayer(clubs Clubs, newClubName, playerName string) (*Player, error) {
	_, p := player(clubs, playerName)
	if p == nil {
		return nil, fmt.Errorf("player %q does not exist", playerName)
	}

	club, err := clubs.Club(newClubName)
	if err != nil {
		return nil, fmt.Errorf("getting new club: %w", err)
	}

	if err := createPlayer(club, p); err != nil {
		return nil, fmt.Errorf("creating player in new club: %w", err)
	}

	if err := removePlayer(clubs, p.Name); err != nil {
		return nil, fmt.Errorf("removing player from old club: %w", err)
	}

	_, p = player(clubs, p.Name)
	if p == nil {
		return nil, fmt.Errorf("could not find player %q after move", p.Name)
	}

	return p, nil
}

func (cs Clubs) UpdatePlayer(p *Player) *Player {
	return nil
}

func updatePlayer(p *Player, up *Player) *Player {
	p.Club = ""

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
func (cs Clubs) BulkUpdatePlayers(ps Players) error {
	return bulkUpdatePlayers(cs, ps)
}

func bulkUpdatePlayers(cs Clubs, ps Players) error {
	// todo: add in moving players
	for _, np := range ps {
		c := maybeMakeClub(cs, np.Club)
		n, p := player(cs, np.Name)
		if n < 0 {
			if err := createPlayer(c, np); err != nil {
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

// move to club file when made“
func maybeMakeClub(cs Clubs, name string) *Club {
	c, ok := cs[name]
	if ok {
		return c
	}

	nc := &Club{Name: name}
	cs[nc.Name] = nc

	return nc
}
