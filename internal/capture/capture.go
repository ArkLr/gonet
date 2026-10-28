package capture

import (
	"net"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

type PacketInfo struct {
	Timestamp       time.Time
	SourceIP        net.IP
	DestinationIP   net.IP
	SourcePort      int
	DestinationPort int
	Protocol        string
	Length          int
}

func Parse(p gopacket.Packet) (PacketInfo, bool) {
	ipLayer := p.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return PacketInfo{}, false
	}
	ip := ipLayer.(*layers.IPv4)

	info := PacketInfo{
		Timestamp:     p.Metadata().Timestamp,
		SourceIP:      ip.SrcIP,
		DestinationIP: ip.DstIP,
		Protocol:      ip.Protocol.String(),
		Length:        len(p.Data()),
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
