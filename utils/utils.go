package utils

import (
	"encoding/csv"
	"errors"
	"strings"

	"github.com/BrianCarducci/DiscordBot/bot_error"
	"github.com/BrianCarducci/DiscordBot/constants"
	"github.com/BrianCarducci/DiscordBot/constants/commands"

	"github.com/bwmarrin/discordgo"
)

var helpStr = help()

func help() string {
	tickmarks := func(s string) string {
		return "`" + s + "`"
	}

	var helpStr string
	validCommands := ""

	keys := []string{}
	for k := range commands.Commands {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return "No commands are available for " + constants.BotName + " yet."
	}
	if len(keys) == 1 {
		return "Usage: " + tickmarks(constants.InvokeStr + " " + keys[0])
	}

	helpStr = "Usage: " + tickmarks(constants.InvokeStr+" [command]") + " where `command` is either "
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

	tokens, err := r.Read()
	if err != nil {
		return []string{}, err
	}

	if tokens[0] != constants.InvokeStr {
		return []string{}, &bot_error.BotError{Message: "", Code: 0}
	}

	if len(tokens) == 1 {
		return []string{}, errors.New("Error: you must call a subcommand.\n" + helpStr)
	}

	return tokens[1:], nil
}

func RunCommand(s *discordgo.Session, msg *discordgo.MessageCreate) (error) {
	tokens, err := tokenize(msg.Content)
	if err != nil {
		return err
	}

	commandStr := tokens[0]
	command, ok := commands.Commands[commandStr]
	if !ok {
		return errors.New("ERROR: " + commandStr + " is not a valid command.\n" + helpStr)
	}
  return command(s, msg, tokens[1:])
}