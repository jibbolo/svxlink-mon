package monitor

import (
	"github.com/jibbolo/svxlink-mon/broker"
	"github.com/jibbolo/svxlink-mon/monitor/reflector"
	"github.com/jibbolo/svxlink-mon/monitor/term"
)

type Monitor struct {
	reflector *reflector.Reflector
	broker    broker.Broker
}

// New create new Monitor instnace
func New(broker broker.Broker) *Monitor {
	return &Monitor{
		reflector.New(),
		broker,
	}

}

func (m *Monitor) Run(quit chan bool) error {

	err := m.broker.Subscribe(m.reflector.Handle)
	if err != nil {
		return err
	}

	go term.Run(m.reflector)

	select {
	case <-quit:
		return nil
	}
}
