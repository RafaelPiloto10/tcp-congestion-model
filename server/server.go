package main

import (
	"log"
	"net"
	"sync"
)

const MAGIC = (512 * (512 + 1)) / 2

type Metrics struct {
	sync.Mutex

	totalErrors uint64

	totalConn   uint64
	droppedConn uint64

	totalPackets   uint64
	droppedPackets uint64
}

func NewMetrics() *Metrics {
	metrics := &Metrics{}
	return metrics
}

func (m *Metrics) AddError(count uint64) {
	m.Lock()
	defer m.Unlock()

	m.totalErrors += count
}

func (m *Metrics) AddDroppedPacket(count uint64) {
	m.Lock()
	defer m.Unlock()

	m.droppedPackets += count
}

func main() {
	metrics := NewMetrics()
	ln, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalf("failed to start tcp server; %v\n", err)
	}

	defer ln.Close()
	log.Printf("tcp server running on 127.0.0.1:8000\n")

	for {

		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept connection; %v\n", err)
			metrics.AddError(1)
		}

		go handleConn(conn, metrics)
	}

}

func handleConn(conn net.Conn, metrics *Metrics) {
	buffer := make([]byte, 512)
	length := 0

	for {
		n, err := conn.Read(buffer)
		length += n

		if err != nil {
			log.Printf("failed to read from buffer; %v\n", err)
			metrics.AddError(1)
			break
		}

		if length >= 512 {
			// TODO: Verify maybe through packet headers that we received a total of 100 packets each sized 512 bytes
			//		with increasing packet ID's and they each contain valid information

			// TODO: May need to increase packet size to force congestion
			break
		}
	}

	if length != 512 {
		metrics.AddDroppedPacket(1)
	}

	if err := conn.Close(); err != nil {
		log.Printf("failed to close connection to %s; %v", conn.RemoteAddr(), err)
		metrics.AddError(1)
	}
}
