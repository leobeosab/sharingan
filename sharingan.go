package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
    "bufio"
    "net"
    "strings"

	"github.com/Ullaakut/nmap"
	"github.com/urfave/cli"
)

func main() {
    DNSBruteForce("yahoo.com", "/home/ae86/ostools/SecLists/Discovery/DNS/namelist.txt")
}

func SetupCLI() {
	sharingan := cli.NewApp()
	sharingan.Name = "Sharingan"
	sharingan.Usage = "Wrapper and analyzer for offensive security recon tools"

	sharingan.Action = func(c *cli.Context) error {
		target := c.Args().Get(0)
		NMAPScan(target)
		return nil
	}

	err := sharingan.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func DNSBruteForce(target string, wordlistPath string) {
    wordlist, err := os.Open(wordlistPath)
    if err != nil {
        log.Fatal(err)
    }

    defer wordlist.Close()

    scanner := bufio.NewScanner(wordlist)
    for scanner.Scan() {
        subdomain := scanner.Text() + "." + target;
        subdomain = strings.Replace(subdomain, " ", "", -1);

        ips, err := ResolveDNS(subdomain)
        if err == nil {
            fmt.Printf("%v", ips)
        } else {
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func ResolveDNS(subdomain string) ([]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    return net.DefaultResolver.LookupHost(ctx, subdomain)
}

func DNSLookup() {

}

func NMAPScan(target string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithContext(ctx),
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
