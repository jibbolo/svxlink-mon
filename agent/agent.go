package agent

import (
	"fmt"
	"io"
	"os"
)

// Agent is responsible for reading file and push its content to
// MQTT server
type Agent struct {
	Path string
	w    io.WriteCloser
}

// New create new Agent instnace
func New(path string, w io.WriteCloser) (*Agent, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid file: %v", err)
	}
	return &Agent{path, w}, nil
}

// Run starts checking the file and push its content forever
func (a *Agent) Run() error {
	f, err := os.Open(a.Path)
	if err != nil {
		return fmt.Errorf("can't open file: %v", err)
	}
	defer f.Close()
	defer a.w.Close()
	return a.cat(f)
}

func (a *Agent) cat(f *os.File) error {
	const NBUF = 512
	var buf [NBUF]byte
	for {
		switch nr, er := f.Read(buf[:]); true {
		case nr < 0:
			return fmt.Errorf("cat: error reading from %s: %s", f.Name(), er.Error())
		case nr == 0: // EOF
			return nil
		case nr > 0:
			if _, err := a.w.Write(buf[0:nr]); err != nil {
				return fmt.Errorf("cat: error writing from %s: %s", f.Name(), err)
			}
		}
	}
}
