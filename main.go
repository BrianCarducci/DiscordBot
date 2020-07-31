package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/BrianCarducci/DiscordBot/services/gunga"
	"github.com/BrianCarducci/DiscordBot/services/weather"
	"github.com/BrianCarducci/DiscordBot/services/odds"

	"github.com/bwmarrin/discordgo"
)

var botName = "JeffBot"
var invokeStr = "!jeff"

var geoLocator = weather.GeoLocator{}

var commands = map[string]func([]string) (string, error) {
	"gunga": gunga.Gunga,
	"weather": geoLocator.GetWeather,
	"odds": odds.PlayOdds
}

var helpStr = help()

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
	geoLocator.Token = gMapsToken

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

	tokens, err := tokenize(m.Content)
	if err != nil {
		return
	}

	if tokens[0] == invokeStr {
		if len(tokens) == 1 {
			s.ChannelMessageSend(m.ChannelID, "Error: you must call a subcommand.\n" + helpStr)
			return
		}

		commandStr := tokens[1]
		command, ok := commands[commandStr]
		if ok == false {
			s.ChannelMessageSend(m.ChannelID, "Error: " + commandStr + " is not a valid command.\n" + helpStr)
			return
		}

		commandArgs := tokens[2:]
		ret, err := command(commandArgs)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: " + commandStr + " failed with message\n" + err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, ret)
	}
}

func help() string {
	tickmarks := func(s string) string {
		return "`" + s + "`"
	}

	var helpStr string
	validCommands := ""

	keys := []string{}
	for k := range commands {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return "No commands are available for " + botName + " yet."
	}
	if len(keys) == 1 {
		return "Usage: " + tickmarks(invokeStr + " " + keys[0])
	}

	helpStr = "Usage: " + tickmarks(invokeStr+" [command]") + " where `command` is either "
	for k := range keys[:len(keys) - 1] {
		validCommands += (tickmarks(keys[k]) + ", ")
	}
	validCommands += ("or " + tickmarks(keys[len(keys)-1]))

	helpStr += validCommands
	return helpStr
}

func tokenize(msg string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(msg))
	r.Comma = ' '

	fields, err := r.Read()
	if err != nil {
		return []string{}, err
	}

	return fields, nil
}
