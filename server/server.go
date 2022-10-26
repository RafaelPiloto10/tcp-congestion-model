package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/RafaelPiloto10/tcp-congestion-model/message"
)

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

func (m *Metrics) AddPacket(count uint64) {
	m.Lock()
	defer m.Unlock()

	m.totalPackets += count
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
	// Force packets to arrive before 3 seconds
	conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	m := message.NewEmptyMessage()
	length := 0

	for {
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)

		for i := 0; i < n; i++ {
			m.Data[length + i] = buffer[i]
		}

		length += n

		metrics.AddPacket(1)

		if err != nil {
			log.Printf("failed to read from buffer; %v\n", err)
			metrics.AddError(1)
			break
		} else if length >= message.BufferLength {
			break
		} else if n != 0 {
			log.Printf("received packet of size %d; total buffer length = %d; buffer = %v\n", n, length, buffer)
		}
	}

	if !m.Checksum() {
		metrics.AddDroppedPacket(1)
		log.Printf("invalid checksum; encountered lost packets; got checksum of %v; expected %v\n", m.GetChecksum(), message.BufferLength*(message.BufferLength+1)/2)
	}

	if err := conn.Close(); err != nil {
		log.Printf("failed to close connection to %s; %v", conn.RemoteAddr(), err)
		metrics.AddError(1)
	}
}
