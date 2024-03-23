package resolver

import (
	"fmt"
	"net"
)

type nameServers []string

var ROOT_NAME_SERVER = "192.5.5.241"

func ResolveDomainName(domainName, rootNameServer string) ([]string, error) {
	dnsMessage, err := queryMessage(rootNameServer, domainName)

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
			if rr.Type == 1 {
				response, err := ResolveDomainName(domainName, rr.Data)

				if err != nil {
					return nil, err
				}

				return response, nil
			}
		}
	}

	if len(dnsMessage.Authority) > 0 {
		var nameServerAddress []string
		for _, rr := range dnsMessage.Authority {
			if rr.Type == 2 {
				nameServerAddress, err = ResolveDomainName(rr.Data, ROOT_NAME_SERVER)
				if err != nil {
					return nil, err
				}
				break
			}
		}

		var response []string = nil
		for _, ip := range nameServerAddress {
			res, err := ResolveDomainName(domainName, ip)
			if err != nil {
				return nil, err
			}
			response = res
		}
		if response != nil {
			return response, nil
		}
	}

	return nil, fmt.Errorf("no answer found")
}

func queryMessage(rootNameServer, domainName string) (*DNSMessage, error) {
	fmt.Printf("Querying %s for %s\n", rootNameServer, domainName)

	ipAddress := net.ParseIP(rootNameServer)
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   ipAddress,
		Port: 53, // Well known port for DNS
	})

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	message, err := GenMessage(domainName)

	if err != nil {
		return nil, err
	}

	_, err = conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("error while sending message: %v", err)
	}

	buffer := make([]byte, 512) // 512 bytes is the maximum size of a DNS message
	_, _, err = conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, fmt.Errorf("error while reading response: %v", err)
	}

	responseMessage, err := DecodeMessage(buffer)
	if err != nil {
		return nil, err
	}

	return &responseMessage, nil
}
