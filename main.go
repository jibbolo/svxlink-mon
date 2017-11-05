package main

import (
	"sync"

	"github.com/alecthomas/kingpin"
	"github.com/jibbolo/svxlink-mon/cmd"
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

	if *agentPath != "" {
		wg.Add(1)
		go cmd.AgentCmd(*agentPath, &wg)
	}

	wg.Wait()
}