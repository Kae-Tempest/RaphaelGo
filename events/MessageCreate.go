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

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	start := time.Now()
	prefix := os.Getenv("PREFIX")
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Contains(m.Content, prefix) {
		m.Content = strings.Replace(m.Content, prefix, "", 1)
	}
	sContent := strings.Split(m.Content, " ")
	cmd := sContent[0]
	sContent = append(sContent[:0], sContent[1:]...)
	switch cmd {
	case "ping":
		Command.Ping(s, m, start)
	case "setup":
		rpg.Setup(s, m, start, sContent)
	case "rename":
		Command.Rename(s, m, start)
	default:
		fmt.Println("Command do not exist !")
	}
}
