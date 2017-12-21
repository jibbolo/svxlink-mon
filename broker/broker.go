package broker

import "io"

// Broker is the message broker for agents and monitor
type Broker interface {
	io.WriteCloser
	Subscribe(fn MessageHandler) error
}

type MessageHandler func([]byte)
