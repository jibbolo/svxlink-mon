package broker

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Broker is the message broker for agents and monitor
type Broker struct {
	Client mqtt.Client
	Topic  string
}

func (b *Broker) Write(p []byte) (int, error) {
	token := b.Client.Publish(b.Topic, 1, false, string(p))
	if token.Wait() && token.Error() != nil {
		return 0, fmt.Errorf("can't publish: %s", token.Error())
	}
	return 0, nil
}

func (b *Broker) Close() error {
	b.Client.Disconnect(250)
	return nil
}

func (b *Broker) Subscribe(topic string, fn MessageHandler) error {
	token := b.Client.Subscribe(topic, 0, b.messageHandler(fn))
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("can't subscribe to %s: %s", topic, token.Error())
	}
	return nil
}

func (b *Broker) messageHandler(fn MessageHandler) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		fn(msg.Payload())
	}
}

type MessageHandler func([]byte)
