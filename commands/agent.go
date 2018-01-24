package cmd

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/jibbolo/svxlink-mon/agent"
)

func AgentCmd(filepath string, follow bool, broker io.WriteCloser, quit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		time.Sleep(time.Second)
		quit <- true
	}()
	var a *agent.Agent
	var err error

	if a, err = agent.New(filepath, broker, follow); err != nil {
		log.Fatalf("Can't init agent: %v", err)
		return
	}
	if err = a.Run(quit); err != nil {
		log.Fatalf("Can't run agent: %v", err)
		return
	}
}