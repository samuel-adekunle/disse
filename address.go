package disse

import (
	"fmt"
	"strings"
)

// Address is a string that identifies a node in the network.
//
// Each node in the network should have a unique address. Nodes can be composed of subnodes, and the uniqueness of addresses is only enforced at the top level.
//
// Subnodes should have addresses that are prefixed with the address of their parent node. This can be done using the NewSubAddress function.
//
// For example, if a node has address "a", and it has a subnode with address "b", then the subnode's full address should be "a.b".
//
// Each subnode address should be unique within the parent node, and this is also enforced by the library.
//
// Additionally, nodes can be nested arbitrarily deep, so a node with address "a.b.c.d" is valid.
type Address string

// GetRoot returns the root address of the given address.
func (a Address) GetRoot() Address {
	return Address(strings.Split(string(a), ".")[0])
}

// GetSubAddress returns a new address which is a valid subnode address.
func (a Address) NewSubAddress(subAddress string) Address {
	return Address(fmt.Sprintf("%s.%s", a, subAddress))
}
