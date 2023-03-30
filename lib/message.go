package lib

type MessageId string
type MessageData interface{}

type Message struct {
	Id   MessageId
	Data MessageData
}

func NewMessage(id MessageId, data MessageData) Message {
	return Message{
		Id:   id,
		Data: data,
	}
}

type MessageTriplet struct {
	Message Message
	From    Address
	To      Address
}
