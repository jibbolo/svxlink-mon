package cmd

import (
	"fmt"
	"sync"

	"github.com/jibbolo/svxlink-mon/broker"
)

func MonitorCmd(broker broker.Broker, wg *sync.WaitGroup) {
	defer wg.Done()

	broker.Subscribe(func(msg []byte) {
		fmt.Printf("--> %s\n", msg)
	})
	select {}

}
