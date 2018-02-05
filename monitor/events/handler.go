package events

import (
	"regexp"
)

const ip = `(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))\:[0-9]+`
const id = `([0-9A-Z]{6})`

// const ID = `(IR0UFQ)`

type Handler struct {
	Comms  chan *Event
	events []*regexp.Regexp
}

var rrr = []*regexp.Regexp{
	regexp.MustCompile(id + ": Login OK from " + ip),
	regexp.MustCompile(id + ": Client " + ip + " disconnected: "),
	regexp.MustCompile(id + ": disconnected: Connection closed by remote peer"),
}

func NewHandler() *Handler {
	return &Handler{
		Comms:  make(chan *Event),
		events: rrr}

}

func (h *Handler) Handle(msg []byte) {
	for _, rgx := range h.events {
		res := rgx.FindSubmatch(msg)
		if res != nil && h.Comms != nil {
			id := string(res[1])
			h.Comms <- &Event{
				ID: id,
			}
			return
		}
	}
}

func (h *Handler) Close() { close(h.Comms) }
