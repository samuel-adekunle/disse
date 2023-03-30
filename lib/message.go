package lib

type MessageId string
type MessageData interface{}

type Message struct {
	Id   MessageId
	Data MessageData
}

type MessageTriplet struct {
	Message Message
	From    Address
	To      Address
}
