package HandleError

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func SendLogError(s *discordgo.Session, m *discordgo.MessageCreate, error string) {
	_, err := s.ChannelMessageSendEmbed("1076795963777220700",
		&discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{Name: m.Author.Username, IconURL: m.Author.Avatar},
			Description: error,
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

func SendDBError(s *discordgo.Session, error string) {
	_, err := s.ChannelMessageSendEmbed("1076795963777220700",
		&discordgo.MessageEmbed{
			Title:       "Database Error",
			Description: error,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "RaphaelGo",
			},
		})
	if err != nil {
		_ = fmt.Errorf("sending Error Message impossible : %s", err)
	}
}
