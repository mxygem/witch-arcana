package main

import (
	"log"

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
	var allClubs bool

	flag.StringVarP(&dataLoc, "data", "d", fileLoc, "location of imported clubs file")
	flag.StringVarP(&clubName, "club", "c", "", "name of player's club")
	flag.StringVarP(&newClubName, "new-club", "m", "", "name of player's new club")
	flag.StringVarP(&name, "name", "n", "", "name of player")
	flag.IntVarP(&level, "level", "l", 0, "player's level")
	flag.IntVarP(&x, "pos-x", "x", 0, "player's x position")
	flag.IntVarP(&y, "pos-y", "y", 0, "player's position")
	flag.StringVar(&csvLoc, "csv", "", "import player data via csv")
	flag.BoolVarP(&allClubs, "all", "a", false, "get all clubs")
	flag.BoolVarP(&shouldLog, "verbose", "v", false, "enable log output")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Println("missing required command & subcommand")
		flag.Usage()
	}

	resource := args[1]
	action := args[0]

	cs := wa.NewClubs(shouldLog)
	if err := cs.LoadData(dataLoc); err != nil {
		log.Fatalf("failed to load data: %v", err)
	}

	switch resource {
	case "club":
		switch action {
		case "get":
			// todo: update printing method
			var d any
			if allClubs {
				d = cs.All()
			} else {
				c, err := cs.Club(clubName)
				if err != nil {
					log.Fatalf("getting club: %v", err)
				}
				d = c
			}
			wa.Print(d)
		case "add":
			c := wa.NewClub(clubName, x, y)

			if err := cs.CreateClub(c); err != nil {
				log.Fatalf("creating club: %v", err)
			}
		case "update":
			c := wa.NewClub(clubName, x, y)
			nc, err := cs.UpdateClub(c)
			if err != nil {
				log.Fatalf("updating club: %v", err)
			}
			wa.Print(nc)
		case "remove":
			if err := cs.RemoveClub(clubName); err != nil {
				log.Fatalf("removing club: %v", err)
			}
		default:
			log.Fatalf("unknown subcommand: %q", action)
		}
	case "player":
		p := wa.NewPlayer(name, clubName, level, x, y)

		switch action {
		case "get":
			gp, err := cs.Player(p.Name)
			if err != nil {
				log.Fatalf("getting player: %v", err)
			}

			print(gp)
		case "add":
			np, err := cs.CreatePlayer(p.Club, p)
			if err != nil {
				log.Fatalf("creating player: %v", err)
			}

			print(np)
		case "update":
			if csvLoc != "" {
				ps, err := wa.ReadCSV(csvLoc)
				if err != nil {
					log.Fatalf("reading players from csv: %v", err)
				}

				if err := cs.BulkUpdatePlayers(ps); err != nil {
					log.Fatalf("bulk updating players: %v", err)
				}
				break
			}

			up, err := cs.UpdatePlayer(p)
			if err != nil {
				log.Fatalf("updating player: %v", err)
			}

			print(up)
		case "move":
			mp, err := cs.MovePlayer(p.Name, newClubName)
			if err != nil {
				log.Fatalf("moving player: %v", err)
			}

			print(mp)
		case "remove":
			if err := cs.RemovePlayer(name); err != nil {
				log.Fatalf("removing player: %v", err)
			}
		default:
			log.Fatalf("unknown subcommand: %q", action)

		}
	default:
		log.Fatalf("unknown resource: %q", resource)
	}

	for _, a := range []string{"add", "update", "remove", "move"} {
		if action != a {
			continue
		}

		save(dataLoc, cs)
	}
}

func print(d any) {
	if err := wa.Print(d); err != nil {
		log.Fatalf("printing: %v", err)
	}
}

func save(loc string, cs *wa.Clubs) {
	if err := wa.Save(loc, cs.All()); err != nil {
		log.Fatalf("saving data: %v", err)
	}
}
