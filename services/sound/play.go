package sound

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/BrianCarducci/DiscordBot/constants"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

var (
	workingDir, _ = os.Getwd()
	soundsDir = filepath.Join(workingDir, "assets", "sounds")
	sounds = map[string]string{
		"name": filepath.Join(soundsDir, "crerb2.pcm"),
	}
	helpString = help()
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
	AudioFrameRate = 48000

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

	EncodeChan = make(chan []int16, 10)
	OutputChan = make(chan []byte, 10)

	OpusEncoder.SetBitrate(AudioBitrate)

	WaitGroup.Add(1)
	go opusEncodeSound()

	WaitGroup.Add(1)
	go playSound(s, m)

	if args[0] == "play" {
		reqSound := args[1]
		soundFilePath, ok := sounds[reqSound]
		if !ok {
			return errors.New("Invalid sound. " + helpString)
		}
		file, err := os.Open(soundFilePath)
		if err != nil {
			fmt.Println("Error opening file :", err)
			return err
		}

		WaitGroup.Add(1)
		go loadSound(file, nil)
	} else {
		pollyAudioStream, err := pollyGetAudioStream(args[1])
		if err != nil {
			return err
		}

		convertSampleRate(pollyAudioStream)
	}

	// wait for above goroutines to finish
	WaitGroup.Wait()

	return nil
}

// TODO: Refactor session creation so it is only created once
func pollyGetAudioStream(message string) (io.ReadCloser, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))		// Create Polly client
	svc := polly.New(sess)

	// Output to PCM using voice Brian
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("pcm"),
		Text: aws.String(message),
		VoiceId: aws.String("Brian"),
		SampleRate: aws.String(constants.PollyAudioSampleRate),
	}

	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		fmt.Println("Got error calling SynthesizeSpeech:")
		return nil, err
	}
	return output.AudioStream, nil
}


// loadSound attempts to load an encoded sound file from disk
// from https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go
func loadSound(reader io.ReadCloser, cmd *exec.Cmd) {
	defer func() {
		close(EncodeChan)
		if cmd != nil {
			cmd.Wait()
		}
		reader.Close()
		WaitGroup.Done()
	}()

	// Create a 16KB input buffer
	fileBuf := bufio.NewReaderSize(reader, 16384)

	// Loop over the file and pass the data to the encoder.
	for {
		buf := make([]int16, AudioFrameSize*AudioChannels)

		fmt.Println("Getting data from audio stream")
		err := binary.Read(fileBuf, binary.LittleEndian, &buf)
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

func convertSampleRate(pollyAudioStream io.ReadCloser) {
	// convert sample rate from Polly to 48KHz using ffmpeg subprocess
	cmd := exec.Command("ffmpeg",
	"-ac", "1",
	"-f", "s16le",
	"-ar", constants.PollyAudioSampleRate,
	"-i", "pipe:",
	"-f", "s16le",
	"-ar", string(AudioFrameRate),
	"-ac", "1",
	"pipe:")

	cmd.Stdin = pollyAudioStream

	ffmpegStdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error assigning ffmpeg stdout pipe:", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting ffmpeg subprocess:", err)
		return
	}

	WaitGroup.Add(1)
	go loadSound(ffmpegStdout, cmd)
}

func opusEncodeSound() {
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
