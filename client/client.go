package main

import (
	"bytes"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")

	if err != nil {
		log.Fatalf("failed to dial to host; %v\n", err)
	}

	defer conn.Close()

	// TODO: Send the server 100 packets, with increasing packet ID
	//		 Consider varying packet size to force congestion?
	buffer := bytes.NewBufferString("Hello world!\n")
	ret, err := conn.Write(buffer.Bytes())

	if err != nil {
		log.Printf("failed to write buffer %b; %v\n", buffer.Bytes(), err)
	}

	if ret < buffer.Len() {
		log.Printf("failed to write complete buffer; got %d; wanted %d\n", ret, buffer.Len())
	}
}
