package cmd

import (
	"log"
	"sync"

	"github.com/jibbolo/svxlink-mon/agent"
)

func AgentCmd(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	var a *agent.Agent
	var err error
	if a, err = agent.New(filepath); err != nil {
		log.Fatalf("Can't init agent: %v", err)
		return
	}
	if err = a.Run(); err != nil {
		log.Fatalf("Can't run agent: %v", err)
		return
	}
}
