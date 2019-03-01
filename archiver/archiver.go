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

	//TODO: Look into using the m->Embeds->url
	url := xurls.Relaxed().FindString(m.Content)
	isRepost, data := fetch(url)

	if isRepost {
		msg := fmt.Sprintf("REPOST ALERT! This url was initially posted by %s on %s channel.", data.User, data.Channel)
		s.ChannelMessageSend(m.ChannelID, msg)
	} else {
		user := m.Author.Username
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			log.WithFields(log.Fields{
				"url": url,
			}).Warnf("Failed to create new entry: %s", err)
		}
		create(url, user, channel.Name)
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

func create(url, user, channel string) {
	return
}
