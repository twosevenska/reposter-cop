package archiver

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"mvdan.cc/xurls/v2"
)

type Metadata struct {
	User    string
	Channel string
}

type Archive struct {
	cache.Cache
}

var archive Archive

// ProcessMessage...
func ProcessMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	//TODO: Look into using the m->Embeds->url
	url := xurls.Relaxed().FindString(m.Content)
	isRepost, data := archive.fetch(url)

	if isRepost {
		msg := fmt.Sprintf("REPOST ALERT! This url was initially posted by %s on %s channel.", data.User, data.Channel)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	user := m.Author.Username
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Warnf("Failed to create new entry: %s", err)
	}

	log.Debugf("Creating new entry for %s", url)
	archive.create(url, user, channel.Name)
}

func Init(expiryTime time.Duration) {
	archive.Cache = *cache.New(expiryTime, 10*time.Minute)
}

func (a Archive) fetch(address string) (bool, Metadata) {
	meta, found := a.Get(address)
	if found {
		return true, meta.(Metadata)
	}
	return false, Metadata{}
}

func (a Archive) create(url, user, channel string) Metadata {
	meta := Metadata{
		User:    user,
		Channel: channel,
	}
	a.Set(url, meta, cache.DefaultExpiration)
	return meta
}
