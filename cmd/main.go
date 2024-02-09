package main

import (
	"fmt"
	"os"

	"github.com/DanielM08/dns-resolver/resolver"
)

func main() {

	domainName := os.Args[1]
	rootNameServer := "192.5.5.241"

	response, err := resolver.ResolveDomainName(domainName, rootNameServer)

	if err != nil {
		fmt.Printf("Error: %s", err)
	} else {
		fmt.Printf("Response: %+v\n", response)
	}
}
