package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/jibbolo/svxlink-mon/broker"
)

func MonitorCmd(wg *sync.WaitGroup) {
	defer wg.Done()

	broker, err := broker.NewAWSBroker("eu-west-1", "", "", "", "mytopic")
	if err != nil {
		log.Fatalf("Can't init monitor: %v", err)
		return
	}

	broker.Subscribe(func(msg []byte) {
		fmt.Printf("--> %s\n", msg)
	})
	select {}

}
