package msg

func FromError(err error) *MessageContainer {
	messages := NewMessageContainer()
	messages.AddErr(err)
	return messages
}
