package events

import (
	"RaphaelGo/Packages/Command"
	"RaphaelGo/Packages/rpg"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
	"time"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate, msg *discordgo.Message) {
	start := time.Now()
	prefix := os.Getenv("PREFIX")
	if msg.Author.ID == s.State.User.ID {
		return
	}
	if strings.Contains(msg.Content, prefix) {
		msg.Content = strings.Replace(msg.Content, prefix, "", 1)
	}
	sContent := strings.Split(msg.Content, " ")
	cmd := sContent[0]
	sContent = append(sContent[:0], sContent[1:]...)
	switch cmd {
	case "ping":
		Command.Ping(s, msg, start)
	case "setup":
		rpg.Setup(s, msg, start, sContent)
	case "rename":
		Command.Rename(s, msg, start)
	default:
		fmt.Println("Command do not exist !")
	}
}
