package console

import (
	"fmt"
	"os"
	"time"

	"github.com/jibbolo/svxlink-mon/broker"
)

type ConsoleBroker struct{}

func New() *ConsoleBroker {
	return &ConsoleBroker{}
}

func (g *ConsoleBroker) Publish(text string) broker.Token {
	_, err := fmt.Fprint(os.Stdout, text)
	return token{err}
}

func (g *ConsoleBroker) Close() {}

type token struct {
	err error
}

func (token) Wait() bool                     { return true }
func (token) WaitTimeout(time.Duration) bool { return true }
func (t token) Error() error                 { return t.err }
