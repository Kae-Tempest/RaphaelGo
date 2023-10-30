package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func sendLogError(s *discordgo.Session, m *discordgo.Message, error string) {
	_, err := s.ChannelMessageSendEmbed("1076795963777220700",
		&discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{Name: m.Author.Username, IconURL: m.Author.Avatar},
			Description: error,
			Timestamp:   time.RFC1123,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    m.Author.ID,
				IconURL: m.Author.Avatar,
			},
		},
	)
	if err != nil {
		_ = fmt.Errorf("sending Error Message impossible : %s", err)
	}
}
