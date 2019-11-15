package nmap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ullaakut/nmap"
)

func FilterHosts(targets *map[string][]string) {
    targetSlice := make([]string, len(*targets))

    i := 0
    for d := range *targets {
        targetSlice[i] = d
        i++
    }

    scanner, err := nmap.NewScanner(
        nmap.WithTargets(targetSlice...),
        nmap.WithPingScan(),
    )

    if err != nil {
		log.Panicf("Unable to create nmap scanner: %v", err)
	}

	result, err := scanner.Run()
	if err != nil {
		log.Panicf("Unable to run nmap scan: %v", err)
	}

    for _, result := range result.Hosts {
        if (result.Status.State == "up") {
            continue
        }

        for _, host := range result.Hostnames {
            delete(*targets, host.Name)
        }
    }
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

	result, err := scanner.Run()
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
