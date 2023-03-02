package lib

// Message type, for all messages in the system
type Message string

// type for triplet of message, from, to
type MessageTriplet struct {
	Message Message
	From    Address
	To      Address
}
