package protocol

import (
	"fmt"
	"net/netip"
	"time"
)

const version = 1

type Type string

const (
	TypeHello  Type = "hello"
	TypeStats  Type = "stats"
	TypePacket Type = "packet"
)

type Message struct {
	Version int        `json:"v"`
	Type    Type       `json:"type"`
	Hello   *HelloMsg  `json:"hello,omitempty"`
	Stats   *StatsMsg  `json:"stats,omitempty"`
	Packet  *PacketMsg `json:"packet,omitempty"`
}

type HelloMsg struct {
	Host      string       `json:"host"`
	Interface string       `json:"iface"`
	LocalIPs  []netip.Addr `json:"local_ips"`
}

// StatsMsg holds counts for one interval, not rates.
type StatsMsg struct {
	Time       time.Time `json:"time"`
	IntervalMs int       `json:"interval_ms"`
	BytesIn    uint64    `json:"bytes_in"`
	BytesOut   uint64    `json:"bytes_out"`
	PacketsIn  uint64    `json:"packets_in"`
	PacketsOut uint64    `json:"packets_out"`
	Dropped    uint64    `json:"dropped"`
}

type PacketMsg struct {
	Time      time.Time  `json:"time"`
	SrcIP     netip.Addr `json:"src_ip"`
	DstIP     netip.Addr `json:"dst_ip"`
	SrcPort   uint16     `json:"src_port"`
	DstPort   uint16     `json:"dst_port"`
	Proto     string     `json:"proto"`
	Length    int        `json:"len"`
	Direction string     `json:"dir"` // "in", "out", "unknown"
}

func (m Message) Validate() error {
	if m.Version != version {
		return fmt.Errorf("version mismatch: got %d, want %d", m.Version, version)
	}
	switch m.Type {
	case TypeStats:
		if m.Stats == nil {
			return fmt.Errorf("%s message has no payload", m.Type)
		}
	case TypeHello:
		if m.Hello == nil {
			return fmt.Errorf("%s message has no payload", m.Type)
		}
	case TypePacket:
		if m.Packet == nil {
			return fmt.Errorf("%s message has no payload", m.Type)
		}

	default:
		return fmt.Errorf("unkown message type %q", m.Type)
	}
	return nil
}
