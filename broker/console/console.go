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
	fmt.Fprint(os.Stdout, text)
	return token{}
}

func (g *ConsoleBroker) Close() {}

type token struct{}

func (token) Wait() bool                     { return true }
func (token) WaitTimeout(time.Duration) bool { return true }
func (token) Error() error                   { return nil }
