package main

import (
	"fmt"
	"log"
	"sort"

	wa "github.com/mxygem/witch-arcana"
	flag "github.com/spf13/pflag"
)

const (
	fileLoc = "clubs.json"
)

var (
	shouldLog bool
)

func main() {
	var clubName, newClubName, name string
	var dataLoc, csvLoc string
	var level, x, y int

	flag.StringVarP(&dataLoc, "data", "d", fileLoc, "location of imported clubs file")
	flag.StringVarP(&clubName, "club", "c", "", "name of player's club")
	flag.StringVarP(&newClubName, "new-club", "m", "", "name of player's new club")
	flag.StringVarP(&name, "name", "n", "", "name of player")
	flag.IntVarP(&level, "level", "l", 0, "player's level")
	flag.IntVarP(&x, "pos-x", "x", 0, "player's x position")
	flag.IntVarP(&y, "pos-y", "y", 0, "player's position")
	flag.StringVar(&csvLoc, "csv", "", "import player data via csv")
	flag.BoolVarP(&shouldLog, "verbose", "v", false, "enable log output")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("missing required command & subcommand")
	}

	cs, err := wa.Open(dataLoc)
	if err != nil {
		log.Fatalf("opening club data: %v", err)
	}

	if shouldLog {
		fmt.Printf("found %d clubs\n", len(cs))
		for k, v := range cs {
			fmt.Printf("club %q with %d players\n", k, len(v.Players))
			for i, p := range v.Players {
				fmt.Printf("\tp %d:%+v\n", i, p)
			}
		}
	}

	switch args[1] {
	case "club":
		c := wa.Club{Name: clubName, Location: wa.Location{X: x, Y: y}}

		if err := manageClub(args[0], cs, c); err != nil {
			log.Fatalf("managing club: %v", err)
		}
	case "player":
		if csvLoc != "" {
			ps, err := wa.ReadCSV(csvLoc)
			if err != nil {
				log.Fatalf("reading players from csv: %v", err)
			}

			if cs == nil {
				log.Println("found nil club data")
				cs = wa.Clubs{}
			}

			if err := cs.BulkUpdatePlayers(ps); err != nil {
				log.Fatalf("bulk updating players: %v", err)
			}
		} else {
			p := &wa.Player{Name: name, Level: level, Location: wa.Location{X: x, Y: y}}

			if err := managePlayer(args[0], cs, clubName, newClubName, p); err != nil {
				log.Fatalf("managing player: %v", err)
			}
		}
	default:
		log.Fatalf("unknown subcommand: %q", args[0])
	}

	if err = wa.Save(dataLoc, cs); err != nil {
		log.Fatalf("saving data: %v", err)
	}
}

func bulkUpdatePlayers(cs wa.Clubs, ps wa.Players) {
	if shouldLog {
		log.Println("-- bulkUpdatePlayers START --")
		log.Printf("cs: %+v\n", cs)
		log.Printf("ps: count: %d\n", len(ps))
	}

	for _, p := range ps {
		c, ok := cs[p.Club]
		if !ok {
			if shouldLog {
				log.Printf("adding new club: %q\n", p.Club)
				log.Printf("and new player: %q\n", p.Name)
			}

			nc := &wa.Club{Name: p.Club, Players: wa.Players{p}}
			cs[nc.Name] = nc

			continue
		}

		if shouldLog {
			log.Printf("looking through club %q's players\n", c.Name)
		}
		var found bool
		for _, cp := range c.Players {
			if p.Name != cp.Name {
				continue
			}

			if shouldLog {
				log.Printf("player match for update: %q\n", p.Name)
			}
			upPlayer(cp, p)
			found = true
			break
		}

		if !found {
			if shouldLog {
				log.Printf("player %q not found, adding\n", p.Name)
			}
			p.Club = ""
			c.Players = append(c.Players, p)
		}
	}

	if shouldLog {
		log.Println("-- bulkUpdatePlayers END --")
	}
}

func upPlayer(p, up *wa.Player) {
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
}

func updatePlayerClubKnown(c *wa.Club, up *wa.Player) {
	if up == nil {
		return
	}

	for _, p := range c.Players {
		if p.Name != up.Name {
			continue
		}

		upPlayer(p, up)

		return
	}
}

func sortedClubIDs(cs wa.Clubs) []string {
	ids := make([]string, 0, len(cs))
	for _, c := range cs {
		ids = append(ids, c.Name)
	}

	return sort.StringSlice(ids)
}

func manageClub(action string, cs wa.Clubs, c wa.Club) error {
	switch action {
	case "get":
		fh, err := getClub(cs, c)
		if err != nil {
			return fmt.Errorf("getting club: %w", err)
		}
		if err = wa.Print(fh); err != nil {
			return fmt.Errorf("printing found club: %w", err)
		}
	case "add":
		if err := addClub(cs, c); err != nil {
			return fmt.Errorf("adding club: %w", err)
		}
	case "update":
		if err := updateClub(cs, c); err != nil {
			return fmt.Errorf("updating club: %w", err)
		}
	case "remove":
		if err := removeClub(cs, c); err != nil {
			return fmt.Errorf("removing club: %w", err)
		}
	default:
		log.Fatalf("unknown subcommand: %q", action)
	}

	return nil
}

func getClub(cd wa.Clubs, c wa.Club) (*wa.Club, error) {
	fh, ok := cd[c.Name]
	if !ok {
		return nil, fmt.Errorf("club %q not found", c.Name)
	}
	if shouldLog {
		log.Printf("found club %q with %d players\n", fh.Name, len(fh.Players))
	}

	return fh, nil
}

func addClub(cs wa.Clubs, c wa.Club) error {
	if _, ok := cs[c.Name]; ok {
		return fmt.Errorf("club %q already exists", c.Name)
	}

	cs[c.Name] = &c

	return nil
}

func updateClub(cs wa.Clubs, c wa.Club) error {
	club, ok := cs[c.Name]
	if !ok {
		return fmt.Errorf("could not find club: %q", c.Name)
	}

	club.Location.X = c.Location.X
	club.Location.Y = c.Location.Y

	return nil
}

func removeClub(cs wa.Clubs, c wa.Club) error {
	if _, ok := cs[c.Name]; !ok {
		return fmt.Errorf("cannot remove nonexistent club: %q", c.Name)
	}

	delete(cs, c.Name)

	return nil
}

func managePlayer(action string, cs wa.Clubs, cName, nCName string, p *wa.Player) error {
	switch action {
	case "get":
		fp, err := getPlayer(cs, cName, p)
		if err != nil {
			return fmt.Errorf("getting player: %w", err)
		}
		if err = wa.Print(fp); err != nil {
			return fmt.Errorf("printing found player: %w", err)
		}
	case "add":
		if err := addPlayer(cs, cName, p); err != nil {
			return fmt.Errorf("adding player: %w", err)
		}
	case "update":
		if err := updatePlayer(cs, cName, p); err != nil {
			return fmt.Errorf("adding player: %w", err)
		}
	case "remove":
		if err := removePlayer(cs, cName, p); err != nil {
			return fmt.Errorf("removing player: %w", err)
		}
	case "move":
		if err := movePlayer(cs, cName, nCName, p); err != nil {
			return fmt.Errorf("moving player: %w", err)
		}
	default:
		return fmt.Errorf("unknown subcommand: %q", action)
	}

	return nil
}

func getPlayer(cd wa.Clubs, cName string, p *wa.Player) (*wa.Player, error) {
	var fc *wa.Club
	if cName != "" {
		ph, err := getClub(cd, wa.Club{Name: cName})
		if err != nil {
			log.Printf("club %q not found\n", cName)
		}
		fc = ph
	}

	var fp *wa.Player
	if fc != nil {
		if _, fp = findPlayer(fc.Players, p.Name); fp != nil {
			return fp, nil
		}

		log.Printf("could not find player %q in club %q, looking through others\n", p.Name, cName)
	}

	for _, c := range cd {
		if _, fp = findPlayer(c.Players, p.Name); fp != nil {
			return fp, nil
		}
	}

	return nil, fmt.Errorf("player %q not found", p.Name)
}

func findPlayer(ps []*wa.Player, pName string) (int, *wa.Player) {
	for i, hp := range ps {
		if hp.Name != pName {
			continue
		}

		return i, hp
	}

	return 0, nil
}

func addPlayer(cs wa.Clubs, cName string, p *wa.Player) error {
	fh, err := getClub(cs, wa.Club{Name: cName})
	if err != nil {
		return fmt.Errorf("club %q not found", cName)
	}

	fp, _ := getPlayer(cs, "", p)
	if fp != nil {
		return fmt.Errorf("player already exists in another club")
	}

	fh.Players = append(fh.Players, p)

	return nil
}

func updatePlayer(cs wa.Clubs, cName string, p *wa.Player) error {
	fp, err := getPlayer(cs, cName, p)
	if err != nil {
		return fmt.Errorf("getting player: %w", err)
	}

	fp.Level = p.Level
	fp.Location.X = p.Location.X
	fp.Location.Y = p.Location.Y
	fp.InHive = p.InHive

	return nil
}

func removePlayer(cs wa.Clubs, cName string, p *wa.Player) error {
	c, err := getClub(cs, wa.Club{Name: cName})
	if err != nil {
		return fmt.Errorf("club %q not found", cName)
	}

	pos, fp := findPlayer(c.Players, p.Name)
	if fp == nil {
		return fmt.Errorf("player not found")
	}

	c.Players = append(c.Players[:pos], c.Players[pos+1:]...)

	return nil
}

func movePlayer(cs wa.Clubs, cName, nCName string, p *wa.Player) error {
	oc, err := getClub(cs, wa.Club{Name: cName})
	if err != nil {
		return fmt.Errorf("original club %q not found", cName)
	}

	nc, err := getClub(cs, wa.Club{Name: nCName})
	if err != nil {
		return fmt.Errorf("new club %q not found", cName)
	}

	_, fp := findPlayer(oc.Players, p.Name)
	if fp == nil {
		return fmt.Errorf("player not found")
	}

	if err = removePlayer(cs, cName, p); err != nil {
		return fmt.Errorf("removing player from original club: %w", err)
	}

	if err = addPlayer(cs, nc.Name, fp); err != nil {
		return fmt.Errorf("adding player to new club: %w", err)
	}

	return nil
}
