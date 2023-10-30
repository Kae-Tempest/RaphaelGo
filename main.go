package main

import (
	"RaphaelGo/events"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	d, err := discordgo.New("Bot " + token)

	d.Identify.Intents = discordgo.IntentsGuildMessages

	d.AddHandler(events.MessageCreate)
	d.AddHandler(events.Ready)

	err = d.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = d.Close()
	if err != nil {
		return
	}
}
