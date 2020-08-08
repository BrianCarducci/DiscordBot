package odds

import (
	"math/rand"
	"errors"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func PlayOdds(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (error) {
	usageError := errors.New("Usage: <\"dare\"> <range> <guess> \nNOTE: <\"dare\"> must be in quotes, <range> and <guess> must be integers")
	msgsend := discordgo.MessageSend{}

	if len(args) != 3 {
		return usageError
	}

	ran, ranErr := strconv.Atoi(args[1])
	guess, guessErr := strconv.Atoi(args[2])

	if ranErr != nil || guessErr != nil {
		return usageError
	}

	if guess < 1 || guess > ran {
		return errors.New("ERROR: Range must be between 1 and " + args[1] + " you sly little dog")
	}

	rngResult := rand.Intn(ran) + 1

	if rngResult == guess {
		msgsend.Content = "LOSER!! xD... Now you must " + args[0]
	} else {
		msgsend.Content = "Ight fine you win dog, cooooool. Bet, let's run it back cuh, you won't."
	}
	s.ChannelMessageSendComplex(m.ChannelID, &msgsend)
	return nil
}