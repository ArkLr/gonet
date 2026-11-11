package main

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/Arklr/gonet/internal/capture"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/afpacket"
	"github.com/gopacket/gopacket/layers"
)

func main() {
	var iface = "wlp1s0"

	localIPs, err := LocalIPs(iface)
	if err != nil {
		log.Fatalf("failed to load local IPs: %v", err)
	}
	fmt.Println("Local IPs:", localIPs)

	localSet := capture.NewLocalSet(localIPs)

	tPacket, err := afpacket.NewTPacket(
		afpacket.OptInterface(iface),
	)
	if err != nil {
		log.Fatalf("Failed to setup afpacket socket: %v", err)
	}
	defer tPacket.Close()

	fmt.Println("Capturing packets via PACKET_MMAP...")

	packets := make(chan capture.PacketInfo, 1000)
	var dropped atomic.Uint64

	// Capture goroutine: read, parse, classify, send without blocking.
	go func() {
		backoff := time.Second

		for {
			handle, err := afpacket.NewTPacket(afpacket.OptInterface(iface))
			if err != nil {
				log.Printf("Open failed: %v, retrying in %s", err, backoff)
				time.Sleep(backoff)
				backoff = min(backoff*2, 30*time.Second)
				continue
			}

			backoff = time.Second
			log.Printf("capture running on %s", iface)
			for {

				data, ci, err := handle.ReadPacketData()
				if err != nil {
					log.Printf("Error reading packet: %v", err)
					break
				}

				gp := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
				gp.Metadata().CaptureInfo = ci

				info, ok := capture.Parse(gp)
				if !ok {
					continue
				}
				info.Direction = capture.GetDirection(info.SourceIP, info.DestinationIP, localSet)

				select {
				case packets <- info:
				default:
					dropped.Add(1)
				}
			}
			handle.Close()
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var stats capture.Stats

	for {
		select {
		case info := <-packets:
			switch info.Direction {
			case capture.DirIn:
				stats.PacketsIn++
				stats.BytesIn += uint64(info.Length)
			case capture.DirOut:
				stats.PacketsOut++
				stats.BytesOut += uint64(info.Length)
			}

		case <-ticker.C:
			fmt.Printf("in: %.1f KB/s (%d pkts)   out: %.1f KB/s (%d pkts)   dropped: %d\n",
				float64(stats.BytesIn)/1024, stats.PacketsIn,
				float64(stats.BytesOut)/1024, stats.PacketsOut,
				dropped.Swap(0),
			)
			stats = capture.Stats{}
		}
	}
}

func LocalIPs(ifaceName string) ([]net.IP, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("interface %q: %w", ifaceName, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("addresses of %q: %w", ifaceName, err)
	}

	var ips []net.IP
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ips = append(ips, ipNet.IP)
	}

	return ips, nil
}
