package models

import "net"

type Clients struct {
	Conn           net.Conn
	CurrentChannel int
	Status         string
}
