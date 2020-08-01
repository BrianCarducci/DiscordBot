package gunga

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var choices = [...]string{"ging", "gung", "gang"}
func Gunga(tokens []string) (*discordgo.MessageSend, error) {
	rand.Seed(time.Now().Unix())

	msg := ""
	for i := 1; i < rand.Intn(50); i++ {
		msg += choices[rand.Intn(len(choices))]
	}

	msgsend := &discordgo.MessageSend{
		Content: msg,
	}
	return msgsend, nil
}
