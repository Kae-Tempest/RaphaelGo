package main

import (
	"RaphaelGo/rpg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type OAuthRes struct {
	Token  string `json:"access_token"`
	Expire int    `json:"expires_in"`
	Type   string `json:"token_type"`
}

type StreamsRes struct {
	Data []struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		UserLogin    string    `json:"user_login"`
		UserName     string    `json:"user_name"`
		GameID       string    `json:"game_id"`
		GameName     string    `json:"game_name"`
		Type         string    `json:"type"`
		Title        string    `json:"title"`
		Tags         []string  `json:"tags"`
		ViewerCount  int       `json:"viewer_count"`
		StartedAt    time.Time `json:"started_at"`
		Language     string    `json:"language"`
		ThumbnailURL string    `json:"thumbnail_url"`
		TagIds       []any     `json:"tag_ids"`
		IsMature     bool      `json:"is_mature"`
	} `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type TwitchUserRes struct {
	Data []struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Type            string    `json:"type"`
		BroadcasterType string    `json:"broadcaster_type"`
		Description     string    `json:"description"`
		ProfileImageURL string    `json:"profile_image_url"`
		OfflineImageURL string    `json:"offline_image_url"`
		ViewCount       int       `json:"view_count"`
		Email           string    `json:"email"`
		CreatedAt       time.Time `json:"created_at"`
	} `json:"data"`
}

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	var StreamOnList []string
	go func() {
		for {
			select {
			case <-ticker.C:
				clientId := os.Getenv("CLIENT_ID")
				bearer := os.Getenv("BEARER")
				clientSecret := os.Getenv("SECRET")
				url := "https://api.twitch.tv/helix/streams?user_login=kaetempest"

				req, err := http.NewRequest("GET", url, nil)
				req.Header.Set("Client-ID", clientId)
				req.Header.Set("Authorization", "Bearer "+bearer)

				client := &http.Client{}
				res, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
				}
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					fmt.Println(err)
				}
				var StreamRes StreamsRes
				json.Unmarshal(resBody, &StreamRes)

				if res.StatusCode == 401 {
					url := fmt.Sprintf(`https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials`, clientId, clientSecret)
					req, err := http.NewRequest("POST", url, nil)
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					client := &http.Client{}
					res, err := client.Do(req)
					if err != nil {
						fmt.Println(err)
					}
					resBody, err := io.ReadAll(res.Body)
					if err != nil {
						fmt.Println(err)
					}
					var data OAuthRes
					json.Unmarshal(resBody, &data)
					os.Setenv("BEARER", data.Token)
					defer res.Body.Close()
				} else if res.StatusCode == 200 {
					if len(StreamRes.Data) > 0 {
						isOn := false
						for u := range StreamOnList {
							if StreamOnList[u] == StreamRes.Data[u].UserLogin {
								isOn = true
							}
						}
						if !isOn {
							StreamOnList = append(StreamOnList, StreamRes.Data[0].UserLogin)
							url := "https://api.twitch.tv/helix/users?login=kaetempest"
							req, err := http.NewRequest("GET", url, nil)
							req.Header.Set("Client-ID", clientId)
							req.Header.Set("Authorization", "Bearer "+bearer)
							client := &http.Client{}
							res, err := client.Do(req)
							if err != nil {
								fmt.Println(err)
							}
							resBody, err := io.ReadAll(res.Body)
							if err != nil {
								fmt.Println(err)
							}
							var TUserRes TwitchUserRes
							json.Unmarshal(resBody, &TUserRes)
							_, err = s.ChannelMessageSendEmbed("1076795963777220700",
								&discordgo.MessageEmbed{
									Author:      &discordgo.MessageEmbedAuthor{Name: StreamRes.Data[0].UserName + " est en stream", IconURL: TUserRes.Data[0].ProfileImageURL},
									Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: TUserRes.Data[0].ProfileImageURL},
									Description: StreamRes.Data[0].Title,
									Fields: []*discordgo.MessageEmbedField{
										{Name: "Viewers", Value: strconv.Itoa(StreamRes.Data[0].ViewerCount), Inline: true},
									},
									Image: &discordgo.MessageEmbedImage{URL: StreamRes.Data[0].ThumbnailURL[:len(StreamRes.Data[0].ThumbnailURL)-20] + "750x450.jpg"},
								})
							if err != nil {
								fmt.Println(err)
							}
						}

					}
				}
				defer res.Body.Close()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

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
		ping(s, m, start)
	case "setup":
		rpg.Setup(s, m, start, sContent)
	case "rename":
		Rename(s, m, start)
	default:
		fmt.Println("Command do not exist !")
	}
}

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
