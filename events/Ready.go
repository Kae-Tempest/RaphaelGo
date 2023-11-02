package events

import (
	"RaphaelGo/Packages/HandleError"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
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

func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	var m *discordgo.MessageCreate
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
				if err != nil {
					errMessage := fmt.Errorf("new Request error for /helix/streams : %s", err).Error()
					HandleError.SendLogError(s, m, errMessage)
				}
				req.Header.Set("Client-ID", clientId)
				req.Header.Set("Authorization", "Bearer "+bearer)

				client := &http.Client{}
				res, err := client.Do(req)
				if err != nil {
					errMessage := fmt.Errorf("new error Do http for /helix/streams: %s", err).Error()
					HandleError.SendLogError(s, m, errMessage)
				}
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					errMessage := fmt.Errorf("new error ResBody for /helix/streanms : %s", err).Error()
					HandleError.SendLogError(s, m, errMessage)
				}
				var StreamRes StreamsRes
				err = json.Unmarshal(resBody, &StreamRes)
				if err != nil {
					errMessage := fmt.Errorf("cannot parse Twitch Stream Response: %s", err).Error()
					HandleError.SendLogError(s, m, errMessage)
				}

				if res.StatusCode == 401 {
					url := fmt.Sprintf(`https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials`, clientId, clientSecret)
					req, err := http.NewRequest("POST", url, nil)
					if err != nil {
						errMessage := fmt.Errorf("new Request error for /oauth/token : %s", err).Error()
						HandleError.SendLogError(s, m, errMessage)
					}
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					client := &http.Client{}
					res, err := client.Do(req)
					if err != nil {
						errMessage := fmt.Errorf("new error Do http for /oauth/token : %s", err).Error()
						HandleError.SendLogError(s, m, errMessage)
					}
					resBody, err := io.ReadAll(res.Body)
					if err != nil {
						errMessage := fmt.Errorf("new error Resbody for /oauth/token : %s", err).Error()
						HandleError.SendLogError(s, m, errMessage)
					}
					var data OAuthRes
					err = json.Unmarshal(resBody, &data)
					if err != nil {
						errMessage := fmt.Errorf("cannot Parse Oauth Twitch Response : %s", err).Error()
						HandleError.SendLogError(s, m, errMessage)
					}
					err = os.Setenv("BEARER", data.Token)
					if err != nil {
						errMessage := fmt.Errorf("error appear when set env key BEARER: %s", err).Error()
						HandleError.SendLogError(s, m, errMessage)
					}
					defer func(Body io.ReadCloser) {
						err := Body.Close()
						if err != nil {
							errMessage := fmt.Errorf("error during closing http connection: %s", err).Error()
							HandleError.SendLogError(s, m, errMessage)
						}
					}(res.Body)
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
							if err != nil {
								errMessage := fmt.Errorf("new Request error for /helix/users : %s", err).Error()
								HandleError.SendLogError(s, m, errMessage)
							}
							req.Header.Set("Client-ID", clientId)
							req.Header.Set("Authorization", "Bearer "+bearer)
							client := &http.Client{}

							res, err := client.Do(req)
							if err != nil {
								errMessage := fmt.Errorf("new error Do http for /helix/users : %s", err).Error()
								HandleError.SendLogError(s, m, errMessage)
							}

							resBody, err := io.ReadAll(res.Body)
							if err != nil {
								errMessage := fmt.Errorf("new error Resbody for /helix/users : %s", err).Error()
								HandleError.SendLogError(s, m, errMessage)
							}

							var TUserRes TwitchUserRes
							err = json.Unmarshal(resBody, &TUserRes)
							if err != nil {
								errMessage := fmt.Errorf("cannot parse Twitch User Response: %s", err).Error()
								HandleError.SendLogError(s, m, errMessage)
							}

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
								errMessage := fmt.Errorf("new error for sending EmbedMessage : %s", err).Error()
								HandleError.SendLogError(s, m, errMessage)
							}
						}
					}
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						errMessage := fmt.Errorf("error during closing http connection: %s", err).Error()
						HandleError.SendLogError(s, m, errMessage)
					}
				}(res.Body)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
