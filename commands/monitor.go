package cmd

import (
	"log"
	"sync"

	"github.com/jibbolo/svxlink-mon/broker"
	"github.com/jibbolo/svxlink-mon/monitor"
)

func MonitorCmd(broker broker.Broker, quit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := monitor.New(broker).Run(quit); err != nil {
		log.Fatalf("Monitor error: %v", err)
	}

}
