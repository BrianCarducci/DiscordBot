package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var botName = "GungaBot"
var choices = [...]string{"ging", "gung", "gang"}
func main() {
	setupBot()
}

func setupBot() {
	envVarName := strings.ToUpper(botName) + "_TOKEN"
	token := strings.TrimSpace(os.Getenv(envVarName))

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Could not instantiate bot. Error: " + err.Error())
		discord.Close()
		os.Exit(1)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: " + err.Error())
		discord.Close()
		os.Exit(1)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println(botName + " is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!gunga" {
		rand.Seed(time.Now().Unix())

		msg := ""
		for i := 1; i < rand.Intn(50); i++ {
			msg += choices[rand.Intn(len(choices))]
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}
