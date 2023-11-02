package events

import (
	"RaphaelGo/Packages/HandleError"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func GuildMemberAdd(s *discordgo.Session, u *discordgo.GuildMemberAdd) {
	var m *discordgo.MessageCreate
	_, err := s.ChannelMessageSendEmbed("1076795963777220700",
		&discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    u.User.GlobalName + " join the server",
				IconURL: u.Avatar,
			},
			Timestamp: time.RFC1123,
		},
	)
	if err != nil {
		errMessage := fmt.Errorf("an error as occured when sending embed: %s", err).Error()
		HandleError.SendLogError(s, m, errMessage)
	}
}

func GuildMemberRemove(s *discordgo.Session, u *discordgo.GuildMemberRemove) {
	var m *discordgo.MessageCreate
	_, err := s.ChannelMessageSendEmbed("1076795963777220700",
		&discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    u.User.GlobalName + " join the server",
				IconURL: u.Avatar,
			},
			Timestamp: time.RFC1123,
		},
	)
	if err != nil {
		errMessage := fmt.Errorf("an error as occured when sending embed: %s", err).Error()
		HandleError.SendLogError(s, m, errMessage)
	}
}
