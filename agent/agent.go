package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/jibbolo/svxlink-mon/broker"
	"github.com/jibbolo/svxlink-mon/broker/google"
)

// Agent is responsible for reading file and push its content to
// MQTT server
type Agent struct {
	Path   string
	Broker broker.Broker
}

// New create new Agent instnace
func New(path string) (*Agent, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid file: %v", err)
	}
	google := google.New()
	return &Agent{path, google}, nil
}

// Run starts checking the file and push its content forever
func (a *Agent) Run() {
	f, err := os.Open(a.Path)
	if err != nil {
		log.Printf("can't open file: %v\n", err)
	}
	defer f.Close()
	defer a.Broker.Close()
	log.Fatal(a.cat(f))
}

func (a *Agent) cat(f *os.File) error {
	const NBUF = 512
	var buf [NBUF]byte
	for {
		switch nr, er := f.Read(buf[:]); true {
		case nr < 0:
			return fmt.Errorf("cat: error reading from %s: %s\n", f.Name(), er.Error())
		case nr == 0: // EOF
			return nil
		case nr > 0:
			token := a.Broker.Publish(string(buf[0:nr]))
			if token.Wait() && token.Error() != nil {
				return fmt.Errorf("cat: error writing from %s: %s\n", f.Name(), token.Error())
			}
		}
	}
}
