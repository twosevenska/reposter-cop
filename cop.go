package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"reposter-cop/archiver"
)

type botConf struct {
	Token      string `envconfig:"token"`
	ExpiryTime int    `envconfig:"expiry_time" default:168`
}

func main() {
	var log = logrus.New()

	var bc botConf
	err := envconfig.Process("discord_cop", &bc)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}
	if bc.Token == "" {
		log.Fatal("Tried loading an empty token")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + bc.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Start our cache
	expiryTime := time.Duration(bc.ExpiryTime) * time.Hour
	archiver.Init(expiryTime)

	// Register the ProcessMessage func as a callback for MessageCreate events.
	dg.AddHandler(archiver.ProcessMessage)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)

	// Cleanly close down the Discord session.
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}
