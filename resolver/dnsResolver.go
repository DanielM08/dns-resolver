package resolver

import (
	"fmt"
	"net"
)

type nameServers []string

func ResolveDomainName(domainName, rootNameServer string) ([]string, error) {

	visitedNS := make(map[string]bool)
	NSInQueue := nameServers{rootNameServer}

	for len(NSInQueue) > 0 {
		ns := NSInQueue[0]
		NSInQueue = NSInQueue[1:]

		response, err := findDomainName(domainName, ns, &NSInQueue, &visitedNS)

		// fmt.Printf("NSInQueue: %+v\n", NSInQueue)

		if err != nil {
			return nil, err
		}

		if response != nil {
			return response, nil
		}
	}

	return nil, fmt.Errorf("No answer found\n")
}

func findDomainName(domainName, rootNameServer string, NSInQueue *nameServers, visitedNS *map[string]bool) ([]string, error) {
	dnsMessage, err := queryMessage(rootNameServer, domainName)

	// fmt.Printf("Response from %s: %+v\n", rootNameServer, dnsMessage)

	// fmt.Printf("dnsMessage: %+v\n", dnsMessage)

	if err != nil {
		return nil, err
	}

	if len(dnsMessage.Answers) > 0 {
		var response []string
		for _, rr := range dnsMessage.Answers {
			response = append(response, rr.Data)
		}
		return response, nil
	}

	if len(dnsMessage.Additional) > 0 {
		for _, rr := range dnsMessage.Additional {
			if rr.Type == 1 && rr.Class == 1 && rr.Length == 4 {
				// fmt.Printf("Adding %s to NSInQueue\n", rr.Data)
				if _, ok := (*visitedNS)[rr.Data]; !ok {
					(*visitedNS)[rr.Data] = true
					*NSInQueue = append(*NSInQueue, rr.Data)
				}
			}
		}
	}

	if len(dnsMessage.Authority) > 0 {
		var nameServerAddress []string
		for _, rr := range dnsMessage.Authority {
			if rr.Type == 2 {
				if _, ok := (*visitedNS)[rr.Data]; !ok {
					(*visitedNS)[rr.Data] = true

					nameServerAddress, err = findDomainName(rr.Data, rootNameServer, NSInQueue, visitedNS)
					if err != nil {
						return nil, err
					}
					break
				}
			}
		}
		if nameServerAddress != nil {
			for _, ip := range nameServerAddress {
				*NSInQueue = append(*NSInQueue, ip)
			}
		}
	}

	return nil, nil
}

func queryMessage(rootNameServer, domainName string) (*DNSMessage, error) {
	fmt.Printf("Querying %s for %s\n", rootNameServer, domainName)

	ipAddress := net.ParseIP(rootNameServer)
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   ipAddress,
		Port: 53, // Well known port for DNS
	})

	if err != nil {
		return nil, fmt.Errorf("Error: %v\n", err)
	}

	defer conn.Close()

	message, err := GenMessage(domainName)

	if err != nil {
		return nil, err
	}

	_, err = conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("Error while sending message: %v\n", err)
	}

	buffer := make([]byte, 512) // 512 bytes is the maximum size of a DNS message
	_, _, err = conn.ReadFromUDP(buffer)
	if err != nil {

		return nil, fmt.Errorf("Error while reading response: %v\n", err)
	}

	responseMessage, err := DecodeMessage(buffer)

	return &responseMessage, nil
}
