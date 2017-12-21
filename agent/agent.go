package agent

import (
	"fmt"
	"io"
	"os"

	"github.com/hpcloud/tail"
)

// Agent is responsible for reading file and push its content to
// MQTT server
type Agent struct {
	Path   string
	w      io.WriteCloser
	follow bool
}

// New create new Agent instnace
func New(path string, w io.WriteCloser, follow bool) (*Agent, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid file: %v", err)
	}
	return &Agent{path, w, follow}, nil
}

// Run starts checking the file and push its content forever
func (a *Agent) Run() error {
	defer a.w.Close()
	return a.tailf(a.follow)
}

func (a *Agent) tailf(follow bool) error {
	t, err := tail.TailFile(a.Path, tail.Config{Follow: follow})
	if err != nil {
		return err
	}
	for line := range t.Lines {
		if _, err := a.w.Write([]byte(line.Text)); err != nil {
			return fmt.Errorf("cat: error writing from %s: %s", a.Path, err)
		}
	}
	return nil
}
