package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/BrianCarducci/DiscordBot/bot_error"
	"github.com/BrianCarducci/DiscordBot/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	setupBot()
}

func getArgs(envNames []string) ([]string) {
	tokenVals := []string{}
	for _, v := range(envNames) {
		val := strings.TrimSpace(os.Getenv(v))
		tokenVals = append(tokenVals, val)
	}
	return tokenVals
}

func setupBot() {
	// Ideally make a map or something for a token's env variable name and value..
	envNames := []string{"DISCORD_TOKEN", "GOOGLE_TOKEN"}
	apiTokens := getArgs(envNames)
	discordToken, gMapsToken := apiTokens[0], apiTokens[1]
	utils.GeoLocator.Token = gMapsToken

	//Exit if one of the needed tokens aren't set
	shouldExit := false
	for i,v := range(apiTokens) {
		if len(v) == 0 {
			fmt.Println(envNames[i] + " environment variable is not set.")
			shouldExit = true
		}
	}
	if shouldExit {
		os.Exit(1)
	}

	discord, err := discordgo.New("Bot " + discordToken)
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
	fmt.Println(utils.BotName + " is now running. Press CTRL-C to exit.")
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

	err := utils.RunCommand(s, m)
	if _, ok := err.(*bot_error.BotError); ok {
		// Might want to check if the code is 0. Hopefully, casting doesn't make a new BotError with code 0..
		return
	} else if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
}
