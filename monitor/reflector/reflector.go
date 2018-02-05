package reflector

import (
	"sort"
	"sync"

	"github.com/jibbolo/svxlink-mon/monitor/parser"
)

// Reflector is the main struct
type Reflector struct {
	parser *parser.Parser
	mutex  *sync.RWMutex
	Links  map[string]*RadioLink
}

// New create new Reflector instance
func New() *Reflector {
	return &Reflector{
		parser: parser.New(),
		mutex:  &sync.RWMutex{},
		Links:  make(map[string]*RadioLink),
	}
}

// Handle handles msg data from MQTT topic
func (r *Reflector) Handle(msg []byte) {
	if event := r.parser.Parse(msg); event != nil {
		r.updateLinkStatus(event)
	}
}

func (r *Reflector) updateLinkStatus(event parser.Event) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	evtID := event.GetID()

	if _, ok := r.Links[evtID]; !ok {
		r.Links[evtID] = &RadioLink{}
	}

	r.Links[evtID].ID = evtID
	r.Links[evtID].TS = event.GetTS()
	if event.GetIP() != "" {
		r.Links[evtID].IP = event.GetIP()
	}

	switch event.(type) {
	case *parser.ClientConnected:
		r.Links[evtID].Status = "idle"
	case *parser.ClientDisconnected:
		r.Links[evtID].Status = "disconnected"
	case *parser.ClientTalkStart:
		r.Links[evtID].Status = "talk"
	case *parser.ClientTalkStop:
		r.Links[evtID].Status = "idle"
	default:
		r.Links[evtID].Status = "unknown"
	}
}

// GetRows returns raw data for UI
func (r *Reflector) GetRows() [][]string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rows := make([][]string, 0)
	for _, l := range r.Links {
		rows = append(rows, []string{
			l.ID, l.IP, l.Status, l.TS,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})
	return rows
}

// RadioLink is the struct for radiolinks
type RadioLink struct {
	ID     string
	IP     string
	TS     string
	Status string
}
