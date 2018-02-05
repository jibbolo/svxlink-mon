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
	defer func() {
		quit <- true
	}()

	err := m.broker.Subscribe(m.reflector.Handle)
	if err != nil {
		return err
	}

	uiQuit := make(chan bool)
	var uiErr error
	go func() {
		uiErr = term.Loop(m.reflector)
		uiQuit <- true
	}()

	select {
	case <-uiQuit:
		return uiErr
	}
}
