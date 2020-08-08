package utils

import (
	"encoding/csv"
	"errors"
	"strings"

	"github.com/BrianCarducci/DiscordBot/bot_error"
	"github.com/BrianCarducci/DiscordBot/services/gunga"
	"github.com/BrianCarducci/DiscordBot/services/weather"
	"github.com/BrianCarducci/DiscordBot/services/odds"
	"github.com/BrianCarducci/DiscordBot/services/m8b"

	"github.com/bwmarrin/discordgo"
)

const BotName = "JeffBot"
const InvokeStr = "!jeff"

var GeoLocator = weather.GeoLocator{}

var commands = map[string]func(*discordgo.Session, *discordgo.MessageCreate, []string) (error) {
	"gunga": gunga.Gunga,
	"weather": GeoLocator.GetWeather,
	"odds": odds.PlayOdds,
	"m8b": m8b.M8b,
}

var helpStr = help()

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
		return "No commands are available for " + BotName + " yet."
	}
	if len(keys) == 1 {
		return "Usage: " + tickmarks(InvokeStr + " " + keys[0])
	}

	helpStr = "Usage: " + tickmarks(InvokeStr+" [command]") + " where `command` is either "
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

	if tokens[0] != InvokeStr {
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
	command, ok := commands[commandStr]
	if !ok {
		return errors.New("ERROR: " + commandStr + " is not a valid command.\n" + helpStr)
	}
  return command(s, msg, tokens[1:])
}