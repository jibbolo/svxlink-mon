package monitor

import (
	"github.com/jibbolo/svxlink-mon/broker"
	"github.com/jibbolo/svxlink-mon/monitor/events"
	"github.com/jibbolo/svxlink-mon/monitor/reflector"
	"github.com/jibbolo/svxlink-mon/monitor/term"
)

type Monitor struct {
	reflector *reflector.Reflector
	handler   *events.Handler
	broker    broker.Broker
}

// New create new Monitor instnace
func New(broker broker.Broker) *Monitor {
	return &Monitor{
		reflector.New(),
		events.NewHandler(),
		broker,
	}

}

func (m *Monitor) Run(quit chan bool) error {

	err := m.broker.Subscribe(m.handler.Handle)
	if err != nil {
		return err
	}

	go term.Run(m.reflector)

	go func() {
		for event := range m.handler.Comms {
			m.reflector.UpdateLinkStatus(event)
		}
	}()

	select {
	case <-quit:
		return nil
	}
}
