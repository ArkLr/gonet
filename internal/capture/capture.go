package capture

import (
	"net"
	"net/netip"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

type Direction string

const (
	DirIn      Direction = "in"
	DirOut     Direction = "out"
	DirUnknown Direction = "unknown"
)

type PacketInfo struct {
	Timestamp       time.Time
	SourceIP        netip.Addr
	DestinationIP   netip.Addr
	SourcePort      int
	DestinationPort int
	Protocol        string
	Length          int
	Direction       Direction
}

type Stats struct {
	BytesIn    uint64
	BytesOut   uint64
	PacketsIn  uint64
	PacketsOut uint64
}

func NewLocalSet(local []net.IP) map[netip.Addr]struct{} {
	set := make(map[netip.Addr]struct{}, len(local))
	for _, ip := range local {
		if a, ok := netip.AddrFromSlice(ip); ok {
			set[a.Unmap()] = struct{}{}
		}
	}
	return set
}

func GetDirection(src, dst netip.Addr, local map[netip.Addr]struct{}) Direction {
	_, srcLocal := local[src.Unmap()]
	_, dstLocal := local[dst.Unmap()]

	switch {
	case dstLocal:
		return DirIn
	case srcLocal:
		return DirOut
	default:
		return DirUnknown
	}
}

func Parse(p gopacket.Packet) (PacketInfo, bool) {
	ipLayer := p.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return PacketInfo{}, false
	}
	ip, ok := ipLayer.(*layers.IPv4)
	if !ok {
		return PacketInfo{}, false
	}

	src, ok := netip.AddrFromSlice(ip.SrcIP)
	if !ok {
		return PacketInfo{}, false
	}
	dst, ok := netip.AddrFromSlice(ip.DstIP)
	if !ok {
		return PacketInfo{}, false
	}

	info := PacketInfo{
		Timestamp:     p.Metadata().Timestamp,
		SourceIP:      src.Unmap(),
		DestinationIP: dst.Unmap(),
		Protocol:      ip.Protocol.String(),
		Length:        p.Metadata().Length,
	}

	if tcpLayer := p.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		info.SourcePort = int(tcp.SrcPort)
		info.DestinationPort = int(tcp.DstPort)
	} else if udpLayer := p.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		info.SourcePort = int(udp.SrcPort)
		info.DestinationPort = int(udp.DstPort)
	}

	return info, true
}
