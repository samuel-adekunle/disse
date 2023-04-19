package disse

import (
	"fmt"

	"github.com/google/uuid"
)

// MessageId is a string that uniquely identifies a message in the network.
//
// UUID is generated using the github.com/google/uuid package.
type MessageId string

// MessageType is a string that identifies a message type and is used to handle messages appropriately.
type MessageType string

// MessageData is the data associated with a message.
type MessageData interface{}

// Message is a message that is sent to a node.
type Message struct {
	Id   MessageId
	Type MessageType
	Data MessageData
}

// String returns a string representation of the message for debugging purposes.
func (m Message) String() string {
	return fmt.Sprintf("%v(%v, %v)", m.Type, m.Id, m.Data)
}

// NewMessage creates a new message with the given messageType and data.
func NewMessage(messageType MessageType, data MessageData) Message {
	return Message{
		Id:   MessageId(uuid.NewString()),
		Type: messageType,
		Data: data,
	}
}

// MessageTriplet is a triplet of a message, the address of the sender and the address of the receiver.
type MessageTriplet struct {
	Message Message
	From    Address
	To      Address
}
