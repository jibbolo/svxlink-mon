package broker

import "time"

// Broker is the message broker for agents and monitor
type Broker interface {
	Publish(string) Token
	Close()
}

// Token is the response you get from Publish, used to
// check errors asynchronously
type Token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}
