package metrics

func messageType(message string, commands []string) string {
	for _, v := range commands {
		if message == "/"+v {
			return v + " command"
		}
	}

	return "text"
}
