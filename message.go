package disse

// MessageId is a string that identifies a message type and is used to handle messages appropriately.
type MessageId string

// MessageData is the data associated with a message.
type MessageData interface{}

// Message is a message that is sent to a node.
type Message struct {
	Id   MessageId
	Data MessageData
}

// NewMessage creates a new message with the given id and data.
func NewMessage(id MessageId, data MessageData) Message {
	return Message{
		Id:   id,
		Data: data,
	}
}

// MessageTriplet is a triplet of a message, the address of the sender and the address of the receiver.
type MessageTriplet struct {
	Message Message
	From    Address
	To      Address
}
