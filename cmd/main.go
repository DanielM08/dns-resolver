package main

import (
	"fmt"

	"github.com/DanielM08/dns-resolver/resolver"
)

func main() {

	domainName := "maxmilhas.com.br"
	rootNameServer := "198.41.0.4"

	// domainName := "codingchallenges.fyi"
	// rootNameServer := "192.36.148.17"

	response, err := resolver.ResolveDomainName(domainName, rootNameServer)

	if err != nil {
		fmt.Printf("Error: %s", err)
	} else {
		fmt.Printf("Response: %+v\n", response)
	}
}
