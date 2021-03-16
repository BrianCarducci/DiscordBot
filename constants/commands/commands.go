package commands

import (
	"github.com/BrianCarducci/DiscordBot/services/gunga"
	"github.com/BrianCarducci/DiscordBot/services/weather"
	"github.com/BrianCarducci/DiscordBot/services/odds"
	"github.com/BrianCarducci/DiscordBot/services/m8b"
	"github.com/BrianCarducci/DiscordBot/services/sound"
	"github.com/bwmarrin/discordgo"
)


var GeoLocator = weather.GeoLocator{}
var Commands = map[string]func(*discordgo.Session, *discordgo.MessageCreate, []string) (error) {
	"gunga": gunga.Gunga,
	"weather": GeoLocator.GetWeather,
	"odds": odds.PlayOdds,
	"m8b": m8b.M8b,
	"play": sound.Play,
	"tts": sound.Play,
}