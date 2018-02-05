package parser

import (
	"regexp"
)

const ip = `(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))\:[0-9]+`
const id = `([0-9A-Z]{6})`

// const ID = `(IR0UFQ)`

type Parser struct {
	events []*regexp.Regexp
}

var rrr = []*regexp.Regexp{
	regexp.MustCompile(id + ": Login OK from " + ip),
	regexp.MustCompile(id + ": Client " + ip + " disconnected: "),
	regexp.MustCompile(id + ": disconnected: Connection closed by remote peer"),
}

func New() *Parser {
	return &Parser{
		events: rrr}

}

func (h *Parser) Parse(msg []byte) *Event {
	for _, rgx := range h.events {
		res := rgx.FindSubmatch(msg)
		if res != nil {
			id := string(res[1])
			return &Event{
				ID: id,
			}
		}
	}
	return nil
}
