package agent

import (
	"fmt"
	"os"

	"github.com/jibbolo/svxlink-mon/broker"
	"github.com/jibbolo/svxlink-mon/broker/aws"
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
	//msgBroker := google.New()
	endpoint := ""
	msgBroker, err := aws.New(endpoint, "eu-west-1", "", "")
	if err != nil {
		return nil, fmt.Errorf("can't start broker: %v", err)
	}
	return &Agent{path, msgBroker}, nil
}

// Run starts checking the file and push its content forever
func (a *Agent) Run() error {
	f, err := os.Open(a.Path)
	if err != nil {
		return fmt.Errorf("can't open file: %v", err)
	}
	defer f.Close()
	defer a.Broker.Close()
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
			token := a.Broker.Publish("mytopic", string(buf[0:nr]))
			if token.Wait() && token.Error() != nil {
				return fmt.Errorf("cat: error writing from %s: %s", f.Name(), token.Error())
			}
		}
	}
}
