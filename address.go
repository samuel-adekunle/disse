package disse

// Address is a string that identifies a node in the network.
//
// Each node in the network should have a unique address.
//
// By convention, subnodes should have addresses that are prefixed with the address of their parent node.
// For example, if a node has address "a", and it has a subnode with address "b", then the subnode's full address should be "a.b".
type Address string
