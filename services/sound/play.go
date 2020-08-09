package sound

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/BrianCarducci/DiscordBot/constants"

	"github.com/bwmarrin/discordgo"
)

var workingDir, _ = os.Getwd()
var soundsDir = filepath.Join(workingDir, "assets", "sounds")
var sounds = map[string]string{
	"name": filepath.Join(soundsDir, "name_jeff.dca"),
}
var helpString = help()

var curSound = ""
var buffer [][]byte 

// Play plays a sound specified by args[0]
func Play(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (error) {
	if len(args) == 0 {
		return errors.New(helpString)
	}

	reqSound := args[0]
	soundFile, ok := sounds[reqSound]
	if !ok {
		return errors.New("Invalid sound. " + helpString)
	}

	if reqSound != curSound {
		err := loadSound(soundFile)
		if err != nil {
			return err
		}
		fmt.Println("PLAYING " + reqSound)

		curSound = reqSound
	}

	err := playSound(s, m)
	if err != nil {
		return err
	}

	//msgsend := discordgo.MessageSend{Content: "play invoked",}

	//s.ChannelMessageSendComplex(m.ChannelID, &msgsend)
	return nil
}

// loadSound attempts to load an encoded sound file from disk
// from https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go
func loadSound(soundFile string) error {
	file, err := os.Open(soundFile)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	buffer = make([][]byte, 0)

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
// from https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go
func playSound(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	// c, err := s.State.Channel(m.ChannelID)
	// if err != nil {
	// 	return errors.New("Could not find channel")
	// }

	// g, err := s.State.Guild(c.GuildID)
	// if err != nil {
	// 	return errors.New("Could not find guild")
	// }
	fmt.Println("Guild ID: " + m.GuildID + "\nChannel ID: " + m.ChannelID)

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(m.GuildID, m.ChannelID, false, false)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

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