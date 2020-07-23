package gunga

import (
	"math/rand"
	"time"
)

var choices = [...]string{"ging", "gung", "gang"}
func Gunga() (string) {
	rand.Seed(time.Now().Unix())

	msg := ""
	for i := 1; i < rand.Intn(50); i++ {
		msg += choices[rand.Intn(len(choices))]
	}
	return msg
}
