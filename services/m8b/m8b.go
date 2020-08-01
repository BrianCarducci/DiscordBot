package m8b

import (
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
	"Not in a milly yearz"
}

func M8b(args []string) (*discordgo.MessageSend, error) {
	question := strings.Join(args, " ")
	answer := answers[rand.Intn(len(answers))]
	msgsend := &discordgo.MessageSend{
		Content: "Question: " + question + "\nAnswer: " + answer,
	}
	return msgsend, nil
}