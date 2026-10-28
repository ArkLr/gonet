package capture

import "time"

type PacketInfo struct {
	timestamp       time.Time
	sourceIP        string
	destinationIP   string
	sourcePort      int
	destinationPort int
	protocol        string
	length          byte
}
