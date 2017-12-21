package broker

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// BasicBroker is the message broker for agents and monitor
type BasicBroker struct {
	Client mqtt.Client
	Topic  string
}

func (b *BasicBroker) Write(p []byte) (int, error) {
	token := b.Client.Publish(b.Topic, 1, false, string(p))
	if token.Wait() && token.Error() != nil {
		return 0, fmt.Errorf("can't publish: %s", token.Error())
	}
	return 0, nil
}

func (b *BasicBroker) Close() error {
	b.Client.Disconnect(250)
	return nil
}

func (b *BasicBroker) Subscribe(fn MessageHandler) error {
	token := b.Client.Subscribe(b.Topic, 1, b.messageHandler(fn))
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("can't subscribe to %s: %s", b.Topic, token.Error())
	}
	return nil
}

func (b *BasicBroker) messageHandler(fn MessageHandler) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		fn(msg.Payload())
	}
}
