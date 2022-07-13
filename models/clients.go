package models

import "net"

type Clients struct {
	conn           net.Conn
	currentChannel int
	status         string
}
