package main

import (
	"fmt"
	"net"

	"github.com/DanielM08/dns-resolver/resolver"
)

func main() {

	message, err := resolver.GenMessage()

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(8, 8, 8, 8),
		Port: 53, // Well known port for DNS
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	defer conn.Close()

	fmt.Println("DNS Server is listening on port: ", conn.LocalAddr().(*net.UDPAddr).Port)

	fmt.Printf("Message sent: %x\n", message)
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error while sending message:", err)
	}

	response := make([]byte, 512) // 512 bytes is the maximum size of a DNS message
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("Error while reading response:", err)
		return
	}

	// Trim the response to actual size
	response = response[:n]

	fmt.Printf("Received response: %x\n", response)
}
