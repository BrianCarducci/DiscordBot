package sound

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	//	"time"

	"github.com/BrianCarducci/DiscordBot/constants"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

var (
	workingDir, _ = os.Getwd()
	soundsDir = filepath.Join(workingDir, "assets", "sounds")
	sounds = map[string]string{
		"name": filepath.Join(soundsDir, "crerb.pcm"),
	}
	helpString = help()

	curSound = ""
)


// These values are taken from https://github.com/bwmarrin/dca/blob/master/cmd/dca/main.go
// Also, the general structure of read -> encode -> output was taken from this repo, but was
// tweaked to work directly with PCM files/streams instead of having to deal with DCA files
var (
	// AudioChannels sets the ops encoder channel value.
	// Must be set to 1 for mono, 2 for stereo
	AudioChannels = 1

	// AudioFrameRate sets the opus encoder Frame Rate value.
	// Must be one of 8000, 12000, 16000, 24000, or 48000.
	// Discord only uses 48000 currently.
	AudioFrameRate = 16000

	// AudioBitrate sets the opus encoder bitrate (quality) value.
	// Must be within 500 to 512000 bits per second are meaningful.
	// Discord only uses 8000 to 128000 and default is 64000.
	AudioBitrate = 64000

	// AudioApplication sets the opus encoder Application value.
	// Must be one of voip, audio, or lowdelay.
	// DCA defaults to audio which is ideal for music.
	// Not sure what Discord uses here, probably voip.
	AudioApplication = "audio"

	// AudioFrameSize sets the opus encoder frame size value.
	// The Frame Size is the length or amount of milliseconds each Opus frame
	// will be.
	// Must be one of 960 (20ms), 1920 (40ms), or 2880 (60ms)
	AudioFrameSize = 960

	// MaxBytes is a calculated value of the largest possible size that an
	// opus frame could be.
	MaxBytes = (AudioFrameSize * AudioChannels) * 2 // max size of opus data

	// OpusEncoder holds an instance of an gopus Encoder
	//OpusEncoder *gopus.Encoder
	//EncoderErr error
	OpusEncoder, EncoderErr = gopus.NewEncoder(AudioFrameRate, AudioChannels, gopus.Audio)

	// WaitGroup is used to wait untill all goroutines have finished.
	WaitGroup sync.WaitGroup

	// Create channels used by the reader/encoder/writer go routines
	EncodeChan = make(chan []int16, 10)
	OutputChan = make(chan []byte, 10)
)

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

	curSound = reqSound

	EncodeChan = make(chan []int16, 10)
	OutputChan = make(chan []byte, 10)

	WaitGroup.Add(1)
	go loadSound(soundFile)

	WaitGroup.Add(1)
	go encodeSound()

	WaitGroup.Add(1)
	go playSound(s, m)

	// wait for above goroutines to finish
	WaitGroup.Wait()

	return nil
}

// loadSound attempts to load an encoded sound file from disk
// from https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go
func loadSound(soundFile string) {
	file, err := os.Open(sounds["name"])
	if err != nil {
		fmt.Println("Error opening pcm file :", err)
	}

	defer func() {
		close(EncodeChan)
		WaitGroup.Done()
		file.Close()
	}()

	// Create a 16KB input buffer
	fileBuf := bufio.NewReaderSize(file, 16384)

	// Loop over the file and pass the data to the encoder.
	for {
		buf := make([]int16, AudioFrameSize*AudioChannels)

		err = binary.Read(fileBuf, binary.LittleEndian, &buf)
		if err == io.EOF {
			// Okay! There's nothing left, time to quit.
			return
		}

		if err == io.ErrUnexpectedEOF {
			// Well there's just a tiny bit left, lets encode it, then quit.
			EncodeChan <- buf
			return
		}

		if err != nil {
			// Oh no, something went wrong!
			fmt.Println("error reading from file,", err)
			return
		}

		// write pcm data to the EncodeChan
		EncodeChan <- buf
	}
}

func encodeSound() {
// encodeSound listens on the EncodeChan and encodes provided PCM16 data
// to opus, then sends the encoded data to the OutputChan

	defer func() {
		close(OutputChan)
		WaitGroup.Done()
	}()

	for {
		pcm, ok := <-EncodeChan
		if !ok {
			// if chan closed, exit
			return
		}

		// try encoding pcm frame with Opus
		opus, err := OpusEncoder.Encode(pcm, AudioFrameSize, MaxBytes)
		if err != nil {
			fmt.Println("Encoding Error:", err)
			return
		}

		// write opus data to OutputChan
		OutputChan <- opus
	}
}

// playSound writer listens on the OutputChan and plays the output to the provided channel.
// from https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go
func playSound(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Could not find channel", err)
	}

	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		fmt.Println("Could not find guild", err)
	}

	var voiceChan string
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			voiceChan = vs.ChannelID
		}
	}
	if len(voiceChan) == 0 {
		fmt.Println("You must be in a voice channel to call this command")
	}

	// If bot is already in the channel, don't do anything
	for _, vs := range g.VoiceStates {
		if vs.UserID == s.State.User.ID {
			return
		}
	}

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(m.GuildID, voiceChan, false, false)
	if err != nil {
		fmt.Println("Couldn't join channel:", err)
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
		WaitGroup.Done()
	}()

	// Sleep for a specified amount of time before playing the sound
	// time.Sleep(32 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	for {
		opus, ok := <-OutputChan
		if !ok {
			// if chan closed, exit
			return
		}
		// Send the buffer data.
		vc.OpusSend <- opus
	}
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
