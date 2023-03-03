package lib

type Message string

type MessageTriplet struct {
	Message Message
	From    Address
	To      Address
}
