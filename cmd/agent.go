package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/jibbolo/svxlink-mon/agent"
)

func AgentCmd(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	var a *agent.Agent
	var err error

	// broker, err := broker.NewAWSBroker("", "eu-west-1", "", "", "mytopic")
	// if err != nil {
	// 	log.Fatalf("Can't init msg broker: %v", err)
	// 	return
	// }
	broker := os.Stdout
	if a, err = agent.New(filepath, broker); err != nil {
		log.Fatalf("Can't init agent: %v", err)
		return
	}
	if err = a.Run(); err != nil {
		log.Fatalf("Can't run agent: %v", err)
		return
	}
}
