package broker

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Broker is the message broker for agents and monitor

type basicBroker struct {
	client mqtt.Client
	topic  string
}

func (b *basicBroker) Write(p []byte) (int, error) {
	token := b.client.Publish(b.topic, 1, false, string(p))
	if token.Wait() && token.Error() != nil {
		return 0, fmt.Errorf("can't publish: %s", token.Error())
	}
	return 0, nil
}

func (b *basicBroker) Close() error {
	b.client.Disconnect(250)
	return nil
}

// Token is the response you get from Publish, used to
// check errors asynchronously
type Token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}
