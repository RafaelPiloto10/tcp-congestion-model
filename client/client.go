package main

import (
	"log"
	"net"

	"github.com/RafaelPiloto10/tcp-congestion-model/message"
)

func RunClient(runs int64) int {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")

	if err != nil {
		log.Fatalf("failed to dial to host; %v\n", err)
	}

	defer conn.Close()

	sentPackets := 0

	for i := 0; int64(i) < runs; i++ {
		buffer := message.NewMessage()
		ret, err := conn.Write(buffer.Data[:])
		if err != nil {
			log.Printf("failed to write buffer of size %d; %v\n", len(buffer.Data), err)
		}

		if ret < len(buffer.Data) {
			log.Printf("failed to write complete buffer; got %d; wanted %d\n", ret, len(buffer.Data))
		} else {
			log.Printf("wrote buffer of size = %d\n", ret)
			sentPackets++
		}
	}

	return sentPackets
}
