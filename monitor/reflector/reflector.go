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

func (r *Reflector) updateLinkStatus(event *parser.Event) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.Links[event.ID]; !ok {
		r.Links[event.ID] = &RadioLink{}
	}

	r.Links[event.ID].ID = event.ID
	r.Links[event.ID].IP = event.IP.String()
	switch event.Action {
	case "connect":
		r.Links[event.ID].Status = "connected"
	case "talk":
		r.Links[event.ID].Status = "talking"
	default:
		r.Links[event.ID].Status = "disconnected"
	}
}

// GetRows returns raw data for UI
func (r *Reflector) GetRows() [][]string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rows := make([][]string, 0)
	for _, l := range r.Links {
		rows = append(rows, []string{
			l.ID, l.IP, l.Status,
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
	Status string
}
