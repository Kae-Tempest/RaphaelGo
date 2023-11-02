package rpg

import (
	"RaphaelGo/Packages/HandleError"
	"RaphaelGo/Packages/Struct"
	"RaphaelGo/Packages/database"
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
	player := database.GetPlayerById(s, p.ID)
	fmt.Println("Player", player)
	//fmt.Println("Player", p)
	fmt.Println("pinged in", timeElapsed)
}

func getCityInfo(ResearchCity string, m *discordgo.MessageCreate, s *discordgo.Session) (Struct.City, error) {
	var selectedCity Struct.City
	CityFind := false
	jsonFile, err := os.Open("./assets/city.json")
	if err != nil {
		return selectedCity, fmt.Errorf("error: %s", err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			errMessage := fmt.Errorf("error during closing json file: %s", err).Error()
			HandleError.SendLogError(s, m, errMessage)
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	var Citys []Struct.City
	err = json.Unmarshal(byteValue, &Citys)
	if err != nil {
		HandleError.SendLogError(s, m, fmt.Errorf("parsing error for city: %s", err).Error())
		return Struct.City{}, err
	}
	for _, city := range Citys {
		if city.Name == ResearchCity && !city.Abandoned {
			selectedCity = city
			CityFind = true
			break
		} else if city.Abandoned {
			_, err := s.ChannelMessageSend(m.ChannelID, "City is abandoned")
			if err != nil {
				errMessage := fmt.Errorf("an Error as occurend when sending message").Error()
				HandleError.SendLogError(s, m, errMessage)
				return Struct.City{}, err
			}
			return selectedCity, errors.New("abandoned")
		}
	}
	if !CityFind {
		_, err := s.ChannelMessageSend(m.ChannelID, "Any city Found")
		if err != nil {
			errMessage := fmt.Errorf("an Error as occurend when sending message").Error()
			HandleError.SendLogError(s, m, errMessage)
			return Struct.City{}, err
		}
		return selectedCity, errors.New("any city found")
	}
	return selectedCity, nil
}

func getJobInfo(ResearchJob string, m *discordgo.MessageCreate, s *discordgo.Session) (Struct.Job, error) {
	var SelectedJob Struct.Job
	JobFind := false
	jsonFile, err := os.Open("./assets/job.json")
	if err != nil {
		return SelectedJob, fmt.Errorf("error: %s", err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			errMessage := fmt.Errorf("error during closing json file: %s", err).Error()
			HandleError.SendLogError(s, m, errMessage)
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	var jobs []Struct.Job
	err = json.Unmarshal(byteValue, &jobs)
	if err != nil {
		HandleError.SendLogError(s, m, fmt.Errorf("parsing error for job: %s", err).Error())
		return Struct.Job{}, err
	}
	for _, job := range jobs {
		if job.Name == ResearchJob {
			SelectedJob = job
			JobFind = true
			break
		}
	}
	if !JobFind {
		_, err := s.ChannelMessageSend(m.ChannelID, "Any Job Find")
		if err != nil {
			errMessage := fmt.Errorf("an Error as occurend when sending message").Error()
			HandleError.SendLogError(s, m, errMessage)
			return Struct.Job{}, err
		}
		return SelectedJob, errors.New("any job find")
	}

	return SelectedJob, nil
}
