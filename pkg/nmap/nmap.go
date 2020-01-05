package nmap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ullaakut/nmap"
	"github.com/leobeosab/sharingan/internal/models"
)

func FilterHosts(targets *[]models.Host) {
	fmt.Printf("\n\nChecking if hosts are up...")

	targetSlice := make([]string, len(*targets))
	// makes deletion trivial
	targetMap := make(map[string]models.Host)

	i := 0
	for _, h := range *targets {
		targetSlice[i] = h.IP
		targetMap[h.IP] = h
		i++
	}

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(targetSlice...),
		nmap.WithPingScan(),
	)

	if err != nil {
		log.Panicf("Unable to create nmap scanner: %v", err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		log.Panicf("Unable to run nmap scan: %v", err)
	}

	// Gather off public internet addresses and discard
	for _, r := range result.Hosts {
		if r.Status.State == "up" {
			continue
		}

		for _, a := range r.Addresses {
			delete(targetMap, a.Addr)
		}
	}
	filtered := make([]models.Host, 0)
	for _, h := range targetMap {
		filtered = append(filtered, h)
	}

	*targets = filtered
}

func Scan(target string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithContext(ctx),
		nmap.WithFastMode(),
	)

	if err != nil {
		log.Fatalf("Unable to create nmap scanner: %v", err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("Unable to run nmap scan: %v", err)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		fmt.Printf("Host %q:\n", host.Addresses[0])

		for _, port := range host.Ports {
			fmt.Printf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
		}
	}

	fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}
