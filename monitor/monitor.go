package main

import (
	"fmt"
	"log"

	"github.com/jibbolo/svxlink-mon/broker"
)

// Thu Sep 28 15:32:09 2017: IW0HKS: Login OK from 195.94.189.122:63358
// Fri Oct  6 20:01:16 2017: IR0UFQ: Client 44.208.124.17:38551 disconnected: Connection closed by remote peer

// Reflector is the main struct
type Reflector struct {
}

// RadioLink is the struct for radiolinks
type RadioLink struct {
	IP string
}

func main() {
	broker, err := broker.NewAWSBroker("", "", "", "", "mytopic")
	if err != nil {
		log.Fatalf("Can't init msg broker: %v", err)
		return
	}
	broker.Subscribe(func(msg []byte) {
		fmt.Printf("--> %s\n", msg)
	})
	select {}
}
