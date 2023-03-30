package disse

import "strings"

// Address is a string that identifies a node in the network.
//
// Nodes can be composed of subnodes and subnodes addresses are concatenated with a dot ('.').
//
// For example 'A.B' is an addresses pointing to node B which is a subnode of node A.
type Address string

// Root returns the address of the root node.
func (a Address) Root() Address {
	return Address(strings.Split(string(a), ".")[0])
}

// SubAddress creates a new address by concatenating the current address with the given address.
//
// The new address is the address of a subnode of the node with the current address.
func (a Address) SubAddress(address Address) Address {
	return a + "." + address
}
