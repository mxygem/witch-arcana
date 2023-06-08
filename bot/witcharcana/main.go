package main

import (
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
)

func main() {
	flag.Parse()

	s, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("creating new session: %v", err)
	}

	cs, err := loadClubs(*DataLocation)
	if err != nil {
		log.Fatalf("could not load data: %v", err)
	}

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

func startUp(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("Bot is up!")
}

func loadClubs(dataLoc string) (*wa.Clubs, error) {
	cs := wa.NewClubs(true)
	if err := cs.LoadData(dataLoc); err != nil {
		return nil, fmt.Errorf("loading data: %w", err)
	}

	return cs, nil
}

func messageHandler(cs *wa.Clubs) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		if m.Content[:4] != "!wat" {
			fmt.Printf("found message start: %q\n", m.Content[:4])
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
