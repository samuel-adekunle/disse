package lib

import "strings"

type Address string

func (a Address) Root() Address {
	return Address(strings.Split(string(a), ".")[0])
}

func (a Address) SubAddress(address Address) Address {
	return a + "." + address
}
