package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup

	fmt.Println("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord")

	for scanner.Scan() {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			checkDomain(domain)
		}(scanner.Text())
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		log.Fatalf("error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {
	fmt.Printf("Checking domain: %v\n", domain)

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up MX records for %s: %v\n", domain, err)
		return
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records for %s: %v\n", domain, err)
		return
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up DMARC records for %s: %v\n", domain, err)
		return
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%v %v %v %v %v %v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
