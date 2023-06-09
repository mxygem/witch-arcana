package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	wa "github.com/mxygem/witch-arcana"
)

const (
	_invalidMsg = "invalid command sent. need action and resource. example: `add player`"
)

func handleMessage(cs *wa.Clubs, m *discordgo.MessageCreate) (any, error) {
	log.Printf("guild: %v channel: %v user: %v\n", m.GuildID, m.ChannelID, m.Author.Username)

	msg := strings.TrimSpace(m.Content[4:])
	log.Printf("trimmed msg: %q\n", msg)

	d := strings.Split(msg, " ")
	if len(d) < 2 {
		return nil, fmt.Errorf(_invalidMsg)
	}

	// todo: update how messages are parsed.
	action := d[0]
	resource := d[1]

	resources := []string{"club", "player"}
	actions := []string{"get", "add", "update", "remove"}
	playerActions := append(actions, "move")

	// once command is valid, set collection as guildID
	cs.SetCollection(m.GuildID)

	switch resource {
	// club
	case resources[0]:
		switch action {
		// get club
		case actions[0]:
			log.Println("get club")
			c, err := cs.Club(d[2])
			if err != nil {
				return nil, fmt.Errorf("getting club: %w", err)
			}

			o, err := wa.PrettyJSON(c)
			if err != nil {
				return nil, fmt.Errorf("formatting data: %w", err)
			}

			return string(o), nil
		// add club
		case actions[1]:
			log.Println("add club")
			var x, y int
			log.Printf("add club arg len: %d\n", len(d))
			if len(d) == 5 {
				xs, err := strconv.Atoi(d[3])
				if err != nil {
					return nil, fmt.Errorf("argument for x: %q is not a valid number", d[3])
				}

				ys, err := strconv.Atoi(d[4])
				if err != nil {
					return nil, fmt.Errorf("argument for y: %q is not a valid number", d[4])
				}
				x = xs
				y = ys
			}

			c := wa.NewClub(d[2], x, y)
			if err := cs.CreateClub(c); err != nil {
				return nil, fmt.Errorf("creating club: %w", err)
			}

			o, err := wa.PrettyJSON(c)
			if err != nil {
				return nil, fmt.Errorf("formatting data: %w", err)
			}

			return string(o), nil
		// update club
		case actions[2]:
			log.Println("update club")
			var x, y int
			if len(d) == 5 {
				xs, err := strconv.Atoi(d[3])
				if err != nil {
					return nil, fmt.Errorf("argument for x: %q is not a valid number", d[3])
				}

				ys, err := strconv.Atoi(d[4])
				if err != nil {
					return nil, fmt.Errorf("argument for y: %q is not a valid number", d[4])
				}
				x = xs
				y = ys
			}
			c := wa.NewClub(d[2], x, y)

			nc, err := cs.UpdateClub(c)
			if err != nil {
				return nil, fmt.Errorf("updating club: %v", err)
			}

			o, err := wa.PrettyJSON(nc)
			if err != nil {
				return nil, fmt.Errorf("formatting data: %w", err)
			}

			return string(o), nil
		// remove club
		case actions[3]:
			fmt.Println("remove club")
			if err := cs.RemoveClub(d[2]); err != nil {
				return nil, fmt.Errorf("removing club: %v", err)
			}
		default:
			return nil, fmt.Errorf("unknown club action %q found. options: %v", action, actions)
		}
	// player
	case resources[1]:
		var clubName string
		if len(d) >= 4 {
			clubName = d[3]
		}

		var level int
		if len(d) >= 5 {
			ls, err := strconv.Atoi(d[4])
			if err != nil {
				return nil, fmt.Errorf("argument for level: %q is not a valid number", d[4])
			}

			level = ls
		}

		var x, y int
		if len(d) >= 7 {
			xs, err := strconv.Atoi(d[5])
			if err != nil {
				return nil, fmt.Errorf("argument for x: %q is not a valid number", d[5])
			}

			ys, err := strconv.Atoi(d[6])
			if err != nil {
				return nil, fmt.Errorf("argument for y: %q is not a valid number", d[6])
			}
			x = xs
			y = ys
		}

		fmt.Printf("player name: %q\n", d[2])
		p := wa.NewPlayer(d[2], clubName, level, x, y)

		switch action {
		// get player
		case playerActions[0]:
			log.Println("get player")
			gp, err := cs.Player(p.Name)
			if err != nil {
				return nil, fmt.Errorf("getting player: %v", err)
			}

			o, err := wa.PrettyJSON(gp)
			if err != nil {
				return nil, fmt.Errorf("formatting data: %w", err)
			}

			return string(o), nil
		// add player
		case playerActions[1]:
			log.Println("add player")
			np, err := cs.CreatePlayer(p.Club, p)
			if err != nil {
				return nil, fmt.Errorf("creating player: %v", err)
			}

			o, err := wa.PrettyJSON(np)
			if err != nil {
				return nil, fmt.Errorf("formatting data: %w", err)
			}

			return string(o), nil
		// update player
		case playerActions[2]:
			log.Println("update player")

			up, err := cs.UpdatePlayer(p)
			if err != nil {
				log.Fatalf("updating player: %v", err)
			}

			o, err := wa.PrettyJSON(up)
			if err != nil {
				return nil, fmt.Errorf("formatting data: %w", err)
			}

			return string(o), nil
		// remove player
		case playerActions[3]:
			log.Println("remove player")

			if err := cs.RemovePlayer(d[3]); err != nil {
				log.Fatalf("removing player: %v", err)
			}
		case playerActions[4]:
			log.Println("move player")
			mp, err := cs.MovePlayer(p.Name, p.Club)
			if err != nil {
				log.Fatalf("moving player: %v", err)
			}

			o, err := wa.PrettyJSON(mp)
			if err != nil {
				return nil, fmt.Errorf("formatting data: %w", err)
			}

			return string(o), nil
		default:
			return nil, fmt.Errorf("unknown player action %q found. options: %v", action, playerActions)
		}
	default:
		return nil, fmt.Errorf("unknown resource %q found. options: %v", resource, resources)
	}

	for _, a := range []string{"add", "update", "remove", "move"} {
		if action != a {
			continue
		}

		if err := wa.Save(cs.DataLocation(), cs.All()); err != nil {
			return nil, fmt.Errorf("saving data: %w", err)
		}
	}

	return nil, nil
}
