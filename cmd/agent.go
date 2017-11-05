package cmd

import (
	"sync"

	"github.com/jibbolo/svxlink-mon/agent"
)

func AgentCmd(filepath string, wg *sync.WaitGroup) error {
	defer wg.Done()
	a, err := agent.New(filepath)
	if err != nil {
		return err
	}
	a.Run()
	return nil
}
