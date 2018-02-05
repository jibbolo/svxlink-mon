package parser

import (
	"regexp"
)

type Event interface {
	GetID() string
	GetIP() string
	GetTS() string
}

type ClientConnected struct{ *baseEvent }
type ClientDisconnected struct{ *baseEvent }
type ClientTalkStart struct{ *baseEvent }
type ClientTalkStop struct{ *baseEvent }

const ip = `(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))\:[0-9]+`
const id = `([0-9A-Z]{6})`

//Sun Sep 24 17:50:18 2017
const ts1 = `([a-zA-z]{3} [a-zA-z]{3} [ 0-9]{2} [\:0-9]{8} [0-9]{4})`

//Sat 03 Feb 2018 04:30:11 PM CET
const ts2 = `([a-zA-z]{3} [0-9]{2} [a-zA-z]{3} [0-9]{4} [\:0-9]{8} [AP]{1}M [A-Z]+)`

const ts = ts2

var rules = []struct {
	rgx *regexp.Regexp
	evt func(res [][]byte) Event
}{
	{regexp.MustCompile(ts + ": " + id + ": Login OK from " + ip), func(res [][]byte) Event {
		return &ClientConnected{
			&baseEvent{
				ts: string(res[1]),
				id: string(res[2]),
				ip: string(res[3]),
			},
		}
	}},
	{regexp.MustCompile(ts + ": " + id + ": Talker start"), func(res [][]byte) Event {
		return &ClientTalkStart{
			&baseEvent{
				ts: string(res[1]),
				id: string(res[2]),
			},
		}
	}},
	{regexp.MustCompile(ts + ": " + id + ": Talker stop"), func(res [][]byte) Event {
		return &ClientTalkStop{
			&baseEvent{
				ts: string(res[1]),
				id: string(res[2]),
			},
		}
	}},
	{regexp.MustCompile(ts + ": " + id + ": Client " + ip + " disconnected: "), func(res [][]byte) Event {
		return &ClientDisconnected{
			&baseEvent{
				ts: string(res[1]),
				id: string(res[2]),
				ip: string(res[3]),
			},
		}
	}},
	{regexp.MustCompile(ts + ": " + id + ": disconnected: Connection closed by remote peer"), func(res [][]byte) Event {
		return &ClientDisconnected{
			&baseEvent{
				ts: string(res[1]),
				id: string(res[2]),
			},
		}
	}},
}

type baseEvent struct {
	ts string
	id string
	ip string
}

func (b *baseEvent) GetID() string {
	return b.id
}

func (b *baseEvent) GetIP() string {
	return b.ip
}

func (b *baseEvent) GetTS() string {
	return b.ts
}
