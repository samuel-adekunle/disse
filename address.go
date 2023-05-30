package disse

// Address is a string that identifies a node in the network.
//
// Each node in the network should have a unique address.
//
// By convention, subnodes should have addresses that are prefixed with the address of their parent node, but this is not enforced by the library.
//
// For example, if a node has address "a", and it has a subnode with address "b", then the subnode's full address should be "a.b".
//
// Additionally, nodes can be nested arbitrarily deep, so a node with address "a.b.c.d" is valid.
type Address string
