package bot_error

type BotError struct {
	Message string
	Code int
}
func (botError *BotError) Error() (string) {
	return botError.Message
}

func New(msg string, code int) (*BotError) {
	return &BotError{msg, code}
}