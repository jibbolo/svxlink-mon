package monitor

import (
	"fmt"

	"github.com/jibbolo/svxlink-mon/broker"
	"github.com/jibbolo/svxlink-mon/monitor/events"
)

type Monitor struct {
	handler *events.Handler
	broker  broker.Broker
}

// New create new Monitor instnace
func New(broker broker.Broker) *Monitor {
	return &Monitor{events.New(), broker}
}

func (m *Monitor) Run(quit chan bool) error {
	err := m.broker.Subscribe(m.handler.Handle)
	if err != nil {
		return err
	}

	go func() {
		for event := range m.handler.Comms {
			fmt.Printf("%+v\n", event)
		}
	}()

	select {
	case <-quit:
		return nil
	}
}

// Thu Sep 28 15:32:09 2017: IW0HKS: Login OK from 195.94.189.122:63358
// Fri Oct  6 20:01:16 2017: IR0UFQ: Client 44.208.124.17:38551 disconnected: Connection closed by remote peer

// Reflector is the main struct
type Reflector struct {
}

// RadioLink is the struct for radiolinks
type RadioLink struct {
	IP string
}
