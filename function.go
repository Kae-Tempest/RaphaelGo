package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"slices"
	"strings"
	"time"
)

func ping(s *discordgo.Session, m *discordgo.MessageCreate, start time.Time) {
	_, err := s.ChannelMessageSend(m.ChannelID, "pong")
	timeElapsed := time.Since(start)
	fmt.Println("pinged in", timeElapsed)
	if err != nil {
		return
	}
}
func Rename(s *discordgo.Session, m *discordgo.MessageCreate, start time.Time) {
	RoleId, RoleName := SearchRole(s, m, "tempest")
	err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, RoleId)
	if err != nil {
		return
	}
	nick := fmt.Sprintf("%s [%s]", m.Author.Username, RoleName)
	errNick := s.GuildMemberNickname(m.GuildID, m.Author.ID, nick)
	if errNick.Error() == "HTTP 403 Forbidden" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error as occurred")
		return
	}
	fmt.Println(RoleId, RoleName)
	timeElapsed := time.Since(start)
	fmt.Println("pinged in", timeElapsed)
}
func SearchRole(s *discordgo.Session, m *discordgo.MessageCreate, researchedRole string) (string, string) {
	var isBot []string
	RoleInfo, _ := s.GuildRoles(m.GuildID)
	GuildMembers, _ := s.GuildMembers(m.GuildID, "", 1000)
	for index, Member := range GuildMembers {
		if Member.User.Bot {
			isBot = append(isBot, Member.Roles[0])
			if index == len(GuildMembers) {
				GuildMembers = slices.Delete(GuildMembers, index-1, index)
			} else {
				GuildMembers = slices.Delete(GuildMembers, index, index+1)
			}
		}
	}
	RoleInfo = append(RoleInfo[:0], RoleInfo[(3):]...)
	for _, Role := range RoleInfo {
		for _, Member := range GuildMembers {
			for _, RoleMember := range Member.Roles {
				if strings.EqualFold(Role.Name, researchedRole) {
					return RoleMember, Role.Name
				}
			}
		}
	}
	return "", ""
}
