package archiver

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"mvdan.cc/xurls/v2"
)

type Metadata struct {
	User    string
	Channel string
}

// ProcessMessage...
func ProcessMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	url := xurls.Relaxed().FindString(m.Content)
	isRepost, data := fetch(url)

	if isRepost {
		msg = fmt.Sprintf("REPOST ALERT! This url was initially posted by %s on %s channel.", data.User, data.Channel)
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	} else {
		err := create(url, user, channel)
		if err != nil {
			log.WithFields(log.Fields{
				"url":     url,
				"user":    user,
				"channel": channel,
			}).Info("Failed to create new entry.")
		}
	}
}

func fetch(address string) (bool, Metadata) {
	if address == "https://imgs.xkcd.com/comics/bad_timing.png" {
		return true, Metadata{
			User:    "garfield",
			Channel: "Techies",
		}
	}
	return false, Metadata{}
}

func create(url, user, channel) error {
	return nil
}
