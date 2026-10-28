package main

import (
	"fmt"
	"log"

	"github.com/Arklr/gonet/internal/capture"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/afpacket"
	"github.com/gopacket/gopacket/layers"
)

func main() {

	tPacket, err := afpacket.NewTPacket(
		afpacket.OptInterface("enp5s0"),
	)
	if err != nil {
		log.Fatalf("Failed to setup afpacket socket: %v", err)
	}
	defer tPacket.Close()

	fmt.Println("Capturing packets via PACKET_MMAP...")

	for {
		data, ci, err := tPacket.ReadPacketData()
		if err != nil {
			log.Printf("Error reading packet: %v", err)
			continue
		}

		gp := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
		gp.Metadata().CaptureInfo = ci
		info, ok := capture.Parse(gp)
		if !ok {
			continue
		}
		fmt.Printf("%s  %s:%d → %s:%d  %s  %d\n",
			info.Timestamp.Format("15:04:05.000"),
			info.SourceIP, info.SourcePort,
			info.DestinationIP, info.DestinationPort,
			info.Protocol,
			info.Length,
		)
	}
}
