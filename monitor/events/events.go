package events

import (
	"net"
)

type Event struct {
	ID     string
	IP     net.IPAddr
	Action string
}
