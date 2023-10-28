package rpg

import (
	"RaphaelGo/database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var db = database.DB()

func Setup(s *discordgo.Session, m *discordgo.MessageCreate, start time.Time, args []string) {
	timeElapsed := time.Since(start)
	if len(args) <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Any args passed")
		return
	}
	JobWanted := args[0]
	job, err := getJobInfo(JobWanted, m, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Any Job find")
		// return error to log-channel
	}
	CityInfo, err := getCityInfo(args[1], m, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Any City find")
		// return error to log-channel
	}

	p := &Player{
		Username: m.Author.GlobalName,
		Level:    1,
		EXP:      0,
		PO:       5,
		Rank:     "G",
		Job:      job.Name,
		City:     CityInfo,
		Attribut: Attribut{
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
		DropRank: DropRank{
			Animals:  strconv.Itoa(rand.Intn(8-1) + 1),
			Plants:   strconv.Itoa(rand.Intn(8-1) + 1),
			Minerals: strconv.Itoa(rand.Intn(8-1) + 1),
			Magic:    strconv.Itoa(rand.Intn(8-1) + 1),
			Special:  strconv.Itoa(rand.Intn(8-1) + 1),
		},
	}
	// insert Player to db
	fmt.Println("Job", job)
	fmt.Println("Player", p)
	fmt.Println("pinged in", timeElapsed)
}

func getCityInfo(ResearchCity string, m *discordgo.MessageCreate, s *discordgo.Session) (City, error) {
	var selectedCity City
	CityFind := false
	jsonFile, err := os.Open("./assets/city.json")
	if err != nil {
		return selectedCity, fmt.Errorf("Error: %s", err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var Citys []City
	json.Unmarshal(byteValue, &Citys)
	for _, city := range Citys {
		if city.Name == ResearchCity && !city.Abandonned {
			selectedCity = city
			CityFind = true
			break
		} else if city.Abandonned {
			s.ChannelMessageSend(m.ChannelID, "City is abandonned")
			return selectedCity, errors.New("Abandonned")
		}
	}
	if !CityFind {
		s.ChannelMessageSend(m.ChannelID, "Any city Found")
		return selectedCity, errors.New("any city found")
	}
	return selectedCity, nil
}

func getJobInfo(ResearchJob string, m *discordgo.MessageCreate, s *discordgo.Session) (Job, error) {
	var SelectedJob Job
	JobFind := false
	jsonFile, err := os.Open("./assets/job.json")
	if err != nil {
		return SelectedJob, fmt.Errorf("Error: %s", err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var jobs []Job
	json.Unmarshal(byteValue, &jobs)
	for _, job := range jobs {
		if job.Name == ResearchJob {
			SelectedJob = job
			JobFind = true
			break
		}
	}
	if !JobFind {
		s.ChannelMessageSend(m.ChannelID, "Any Job Find")
		return SelectedJob, errors.New("any job find")
	}

	return SelectedJob, nil
}
