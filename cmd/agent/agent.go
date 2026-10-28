package main

import (
	"fmt"
	"log"

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

	fmt.Println("Capturing packets via PACKET_MAP...")

	for {
		data, _, err := tPacket.ReadPacketData()
		if err != nil {
			log.Printf("Error reading packet: %v", err)
			continue
		}

		gp := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
		fmt.Print("packets: \n", gp)
	}
}
