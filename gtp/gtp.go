package gtp

func Completions(sender string, msg string) (string, error) {
	//reply, err := Minimax_conversation(sender, msg)
	reply, err := Xinghuo_conversation(sender, msg)

	return reply, err
}
