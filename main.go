package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/briancarducci/DiscordBot/services/gunga"
)

var botName = "JeffBot"
var invokeStr = "!jeff"
var commands = map[string] map[string]interface{} {
	"gunga": {
		"nArgs": 0,
		"func": gunga.Gunga,
	},
}
var helpStr = help()

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

	tokens, err = tokenize(m.Content)
	if err != nil {
		return
	}

	if tokens[0] == invokeStr {
		commandStr := tokens[1]
		command, ok := commands[commandStr]
		if ok != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: " + commandStr " is not a valid command. " + helpStr)
			return
		}

		nArgs, aErr, cFunc, cFuncErr := command["nArgs"], command["func"]
		if aErr != nil {
			fmt.Println("Error: Command " + commandStr + " doesn't have an 'nArgs' field.")
			s.ChannelMessageSend(m.ChannelID, "Error. Check bot logs for details.")
			return
		}
		if cFuncErr != nil {
			fmt.Println("Error: Command " + commandStr + " doesn't have an 'func' field.")
			s.ChannelMessageSend(m.ChannelID, "Error. Check bot logs for details.")
			return
		}

		commandArgs := tokens[2:]
		nArgsEqual := len(commandArgs) == nArgs
		var ret string
		if nArgsEqual && len(commandArgs) == 0 {
			ret = commandFunc()
		}
		else if nArgsEqual && len(commandArgs) > 0 {
			ret = commandFunc(commandArgs)
		}
		else {
			fmt.Println("Error: " + commandStr + " takes " + nArgs + " and was called with " + len(commandArgs) + " args but something went wrong.\n")
			s.ChannelMessageSend(m.ChannelID, "Error. Check bot logs for details.")
			return
		}

		s.ChannelMessageSend(m.ChannelID, ret)
	}
}

func help() (string) {
	tickmarks := func (s string) (string) {
		return "`" + s + "`"
	}

	helpStr := "Please enter " + tickmarks(invokeStr + "[command]") + " where `command` is either "
	validCommands := ""
	for k := range commands[:len(commands)-1] {
		validCommands += (tickmarks(k) + ", ")
	}
	validCommands += ("or " + tickmarks(commands[len(commands)]))

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
