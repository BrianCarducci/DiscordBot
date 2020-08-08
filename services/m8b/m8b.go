package m8b

import (
	"errors"
	"math/rand"
	"strings"
	"github.com/bwmarrin/discordgo"
)

var answers = []string{
	"It is certain.",
	"Not even close baaabyyy.",
	"Naaaaah",
	"idk bruh",
	"Yuh",
	"Probs",
	"Not in a milly yearz",
}

func M8b(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (error) {
	if len(args) == 0 {
		return errors.New("Provide a question, yo")
	}
	question := strings.Join(args, " ")
	answer := answers[rand.Intn(len(answers))]
	msgsend := &discordgo.MessageSend{
		Content: "Question: " + question + "\nAnswer: " + answer,
	}
	s.ChannelMessageSendComplex(m.ChannelID, msgsend)
	return nil
}