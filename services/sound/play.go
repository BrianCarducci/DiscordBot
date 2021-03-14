package sound

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	//	"time"

	"github.com/BrianCarducci/DiscordBot/constants"

	"github.com/bwmarrin/discordgo"
)

var workingDir, _ = os.Getwd()
var soundsDir = filepath.Join(workingDir, "assets", "sounds")
var sounds = map[string]string{
	"name": filepath.Join(soundsDir, "crerb.pcm"),
}
var helpString = help()

var curSound = ""
var buffer []byte

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

		curSound = reqSound
	}

	err := playSound(s, m)
	if err != nil {
		return err
	}

	return nil
}

// loadSound attempts to load an encoded sound file from disk
// from https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go
func loadSound(soundFile string) error {

	fileDetails, err := os.Stat(filepath.Join(soundsDir, "crerb.pcm"))
	if err != nil {
		fmt.Println("Error opening pcm file :", err)
		return err
	}
	fileSize := fileDetails.Size()

	buffer = make([]byte, fileSize)

	file, err := os.Open(filepath.Join(soundsDir, "crerb.pcm"))
	if err != nil {
		fmt.Println("Error opening pcm file :", err)
		return err
	}

	err = binary.Read(file, binary.LittleEndian, &buffer)
	if err != nil {
		fmt.Println("Error reading pcm file :", err)
		return err
	}

	return nil

	// for {
	// 	// Read opus frame length from dca file.
	// 	err = binary.Read(file, binary.LittleEndian, &opuslen)

	// 	// If this is the end of the file, just return.
	// 	if err == io.EOF || err == io.ErrUnexpectedEOF {
	// 		err := file.Close()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	}

	// 	if err != nil {
	// 		fmt.Println("Error reading from dca file :", err)
	// 		return err
	// 	}

	// 	// Read encoded pcm from dca file.
	// 	InBuf := make([]byte, opuslen)
	// 	err = binary.Read(file, binary.LittleEndian, &InBuf)

	// 	// Should not be any end of file errors
	// 	if err != nil {
	// 		fmt.Println("Error reading from dca file :", err)
	// 		return err
	// 	}

	// 	// Append encoded pcm data to the buffer.
	// 	buffer = append(buffer, InBuf)
	// }
}

// playSound plays the current buffer to the provided channel.
// from https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go
func playSound(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return errors.New("Could not find channel")
	}

	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return errors.New("Could not find guild")
	}

	var voiceChan string
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			voiceChan = vs.ChannelID
		}
	}
	if len(voiceChan) == 0 {
		return errors.New("You must be in a voice channel to call this command")
	}

	// If bot is already in the channel, don't do anything
	for _, vs := range g.VoiceStates {
		if vs.UserID == s.State.User.ID {
			return nil
		}
	}

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(m.GuildID, voiceChan, false, false)
	if err != nil {
		return err
	}

	// The thought here was to wait until 5 seconds elapsed and if bot was still in the channel, exit.
	// Doesn't seem to be working though
	//
	// t0 := time.Now()
	// inChannel := false
	// for !inChannel && time.Since(t0).Seconds() < 5 {
	// 	if time.Since(t0).Seconds() > 5 {
	// 		return errors.New("Timed out bruh")
	// 	}
	// 	for _, vs := range g.VoiceStates {
	// 		if vs.UserID == s.State.User.ID {
	// 			inChannel = true
	// 			continue
	// 		}
	// 	}
	// }

	defer func() {
		// Stop speaking
		vc.Speaking(false)

		// Sleep for a specified amount of time before ending.
		// time.Sleep(32 * time.Millisecond)

		// Disconnect from the provided voice channel.
		vc.Disconnect()
	}()

	// Sleep for a specified amount of time before playing the sound
	// time.Sleep(32 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	vc.OpusSend <- buffer
	

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
