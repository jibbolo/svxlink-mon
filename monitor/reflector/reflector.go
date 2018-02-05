package reflector

import (
	"sort"
	"sync"

	"github.com/jibbolo/svxlink-mon/monitor/events"
)

// Reflector is the main struct
type Reflector struct {
	mutex *sync.RWMutex
	Links map[string]*RadioLink
}

func New() *Reflector {
	return &Reflector{
		mutex: &sync.RWMutex{},
		Links: make(map[string]*RadioLink),
	}
}

func (r *Reflector) UpdateLinkStatus(event *events.Event) {
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
