package message

type MessageError struct {
	Msg string
}

func (receiver MessageError) Error() string {
	return receiver.Msg
}

func NewMessageError(Msg string) *MessageError {
	return &MessageError{Msg: Msg}
}
