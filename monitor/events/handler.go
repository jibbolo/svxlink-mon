package events

import (
	"fmt"
	"regexp"
	"sync"
)

// import (
// 	"net"
// 	"regexp"
// 	"time"
// )

// type ParserFunc func(string) *Message

// type Message struct {
// 	IP net.IPAddr
// 	TS time.Time
// 	ID string
// }

// type MessageType struct {
// 	regex  *regexp.Regexp
// 	parser ParserFunc
// }

// type Handler struct {
// 	events []*MessageType
// }

// func (h *Handler) Register(rgx *regexp.Regexp, parser ParserFunc) {
// 	h.events = append(h.events, &MessageType{
// 		regex:  rgx,
// 		parser: parser,
// 	})
// }

const IP = `(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))\:[0-9]+`

const ID = `([0-9A-Z]{6})`

// const ID = `(IR0UFQ)`

var clientConnected = regexp.MustCompile("Client " + IP + " connected")
var loginOK = regexp.MustCompile(ID + ": Login OK from " + IP)
var clientDisconnect = regexp.MustCompile(ID + ": Client " + IP + " disconnected: ")
var clientDisconnect2 = regexp.MustCompile(ID + ": disconnected: Connection closed by remote peer")

type Handler struct{}

var lock = &sync.Mutex{}
var ids = make(map[string][]string)

func (h *Handler) Parse(msg []byte) {
	lock.Lock()
	defer lock.Unlock()
	var id string

	res := clientDisconnect.FindSubmatch(msg)
	if res != nil {
		id = string(res[1])
		ids[id] = append(ids[id], "disconnect")
		return
	}

	res = clientDisconnect2.FindSubmatch(msg)
	if res != nil {
		id = string(res[1])
		ids[id] = append(ids[id], "disconnect")
		return
	}

	res = loginOK.FindSubmatch(msg)
	if res != nil {
		id = string(res[1])
		ids[id] = append(ids[id], "connect")
		return
	}

}
func (h *Handler) String() string {
	return fmt.Sprintf("%s", ids)
}

var DefaultHandler = &Handler{}
