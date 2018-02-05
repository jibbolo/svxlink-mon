package main

import (
	"sync"

	"github.com/alecthomas/kingpin"
	"github.com/jibbolo/svxlink-mon/broker"
	"github.com/jibbolo/svxlink-mon/commands"
)

var build = "dev"

var (
	debug = kingpin.Flag("debug", "Enable debug mode.").
		Bool()
	monitor = kingpin.Flag("monitor", "Enable monitor mode.").
		Short('m').
		Bool()
	agentPath = kingpin.Flag("agent", "Enable agent mode and requires a file to read from.").
			Short('a').
			ExistingFile()
)

func main() {
	kingpin.Version(build)
	kingpin.Parse()

	var wg sync.WaitGroup

	broker := broker.NewLoopback()

	quit := make(chan bool)

	if *agentPath != "" {
		wg.Add(1)
		go cmd.AgentCmd(*agentPath, broker, quit, &wg)
	}

	if *monitor {
		wg.Add(1)
		go cmd.MonitorCmd(broker, quit, &wg)
	}

	wg.Wait()
}
