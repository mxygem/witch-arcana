package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	wa "github.com/mxygem/witch-arcana"
)

var (
	BotToken     = flag.String("token", "", "Bot access token")
	DataLocation = flag.String("data", "", "Location of data")
	DBConn       = flag.String("dbconn", "", "database location")
	DBName       = flag.String("dbname", "", "database name")
	Debug        = flag.Bool("debug", false, "enable verbose logging")
)

func main() {
	flag.Parse()

	s, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("creating new session: %v", err)
	}

	// cs, err := loadClubs(*DataLocation)
	// if err != nil {
	// 	log.Fatalf("could not load data: %v", err)
	// }

	ctx := context.Background()
	db := wa.NewDB(ctx, *DBConn, *DBName)
	if err := db.Connect(); err != nil {
		log.Fatalf("connecting to db: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db unreachable: %v", err)
	}

	cs := wa.NewClubs(db, *Debug)

	s.AddHandler(startUp)
	s.AddHandler(messageHandler(cs))

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down")
}

// func loadClubs(dataLoc string) (*wa.Clubs, error) {
// 	cs := wa.NewClubs(true)
// 	if err := cs.LoadData(dataLoc); err != nil {
// 		return nil, fmt.Errorf("loading data: %w", err)
// 	}

// 	return cs, nil
// }

func startUp(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("Bot is up!")
}

func messageHandler(cs *wa.Clubs) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// filter bots or messages not intended for this one.
		if m.Author.Bot || m.Content[:4] != "!wat" {
			return
		}

		d, err := handleMessage(cs, m.Content)
		if err != nil {
			msg := fmt.Sprintf("bad request: %v", err)
			_, err := s.ChannelMessageSend(m.ChannelID, msg)
			if err != nil {
				log.Printf("could not send error message: %v", err)
			}
		}

		if d != nil {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprint(d))
			if err != nil {
				log.Printf("could not send data message: %v", err)
			}
		}
	}
}
