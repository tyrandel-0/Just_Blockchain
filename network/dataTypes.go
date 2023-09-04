package network

import "net"

type Package struct {
	Option int
	Data   string
}

type Listener net.Listener
type Conn net.Conn
