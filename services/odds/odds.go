package odds

import (
	"strings"
	"math/rand"
	"errors"
	"strconv"
)

func PlayOdds(args []string) (string, error) {
	usageError := errors.New("Usage: <\"dare\"> <range> <guess> \nNOTE: <\"dare\"> must be in quotes, <range> and <guess> must be integers")
	if len(args) != 3 {
		return "", usageError
	}

	ran, err := strconv.Atoi(args[1])
	if err != nil {
		return "", usageError
	}
	guess, err := strconv.Atoi(args[2])
	if err != nil {
		return "", usageError
	}

	rngResult := rand.Intn(ran) + 1

	if rngResult == guess {
		return "LOSER!! xD... Now you must " + args[0], nil
	} else {
		return "Ight fine you win dog, cooooool. Bet, let's run it back cuh, you won't.", nil
	}
}