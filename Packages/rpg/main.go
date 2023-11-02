package rpg

import (
	"RaphaelGo/Packages/HandleError"
	"RaphaelGo/Packages/Struct"
	"RaphaelGo/Packages/database"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"time"
)

func Setup(s *discordgo.Session, m *discordgo.MessageCreate, start time.Time, args []string) {
	timeElapsed := time.Since(start)
	if len(args) <= 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Any args passed")
		if err != nil {
			errMessage := fmt.Errorf("an Error as occurend when sending message").Error()
			HandleError.SendLogError(s, m, errMessage)
		}
		return
	}
	job := database.GetJobByName(s, args[0])
	CityInfo := database.GetCityByName(s, args[1])

	p := Struct.Player{
		Username: m.Author.GlobalName,
		Level:    1,
		EXP:      0,
		PO:       5,
		Rank:     "G",
		Job:      job,
		City:     CityInfo,
		Attribute: Struct.Attribute{
			HP:           25,
			MP:           10,
			Strength:     job.Strength,
			Agility:      job.Agility,
			Intelligence: job.Intelligence,
			Spirit:       job.Spirit,
			Constitution: job.Constitution,
			Speed:        job.Speed,
			Luck:         rand.Intn(5-1) + 1,
		},
		DropRank: Struct.DropRank{
			Animals:  strconv.Itoa(rand.Intn(8-1) + 1),
			Plants:   strconv.Itoa(rand.Intn(8-1) + 1),
			Minerals: strconv.Itoa(rand.Intn(8-1) + 1),
			Magic:    strconv.Itoa(rand.Intn(8-1) + 1),
			Special:  strconv.Itoa(rand.Intn(8-1) + 1),
		},
	}
	// insert Player to db
	database.CreatePlayer(s, p)
	player := database.GetPlayerByName(s, p.Username)
	fmt.Println("Player", player)
	//fmt.Println("Player", p)
	fmt.Println("pinged in", timeElapsed)
}
