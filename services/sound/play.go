package sound

import (
	"fmt"
	"errors"
	"os"
	"path/filepath"

	"github.com/BrianCarducci/DiscordBot/constants"

	"github.com/bwmarrin/discordgo"
)

var workingDir, _ = os.Getwd()
var soundsDir = filepath.Join(workingDir, "assets", "sounds")
var sounds = map[string]string{
	"name": filepath.Join(soundsDir, "name_jeff.mp3"),
}
var helpString = help()

func Play(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (error) {
	if len(args) == 0 {
		return errors.New(helpString)
	}

	//msgsend := discordgo.MessageSend{Content: "play invoked",}

	//s.ChannelMessageSendComplex(m.ChannelID, &msgsend)
	return nil
}

func help() string {
	tickmarks := func(s string) string {
		return "`" + s + "`"
	}

	var helpStr string
	validCommands := ""

	keys := []string{}
	for k := range sounds {
		keys = append(keys, k)
	}
	
	if len(keys) == 1 {
		return "Usage: " + tickmarks(constants.InvokeStr + " play " + keys[0])
	}

	helpStr = "Usage: " + tickmarks(constants.InvokeStr + " play " + "[sound]") + " where `[sound]` is either "
	for k := range keys[:len(keys) - 1] {
		validCommands += (tickmarks(keys[k]) + ", ")
	}
	validCommands += ("or " + tickmarks(keys[len(keys)-1]))

	helpStr += validCommands
	return helpStr
}