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

    "text/tabwriter"

	"github.com/Ullaakut/nmap"
	"github.com/urfave/cli"
)

func main() {
    SetupCLI()
}

func SetupCLI() {

    var dnsWordlist string
    var target string

	sharingan := cli.NewApp()
	sharingan.Name = "Sharingan"
	sharingan.Usage = "Wrapper and analyzer for offensive security recon tools"

    sharingan.Flags = []cli.Flag {
        cli.StringFlag{
            Name: "dns-wordlist",
            Value: "",
            Usage: "Wordlist for DNS bruteforcing",
            Destination: &dnsWordlist,
        },
        cli.StringFlag{
            Name: "target",
            Value: "",
            Usage: "Target domain",
            Destination: &target,
        },
    }

	sharingan.Action = func(c *cli.Context) error {
        RunDNSRecon(target, dnsWordlist)
		return nil
	}

	err := sharingan.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func RunDNSRecon(target string, wordlistPath string) {

    if (target == "") {
        log.Fatal("Target needs to be defined")
    }

    if (wordlistPath == "") {
        log.Fatal("DNS Wordlist needs to be defined")
    }

    w := new(tabwriter.Writer)
    w.Init(os.Stdout, 8, 8, 0, '\t', 0)

    validSubdomains := DNSBruteForce(target, wordlistPath)

    fmt.Fprintf(w, "\n%s\t%s\t", "Subdomain Address ", "| IP List")
    fmt.Fprintf(w, "\n%s\t%s\t", "----------------- ", "| -------")
    for subdomain, ips := range validSubdomains {
        fmt.Fprintf(w, "\n%s \t| %v\t", subdomain, ips)
    }

    w.Flush()
}

func DNSBruteForce(target string, wordlistPath string) (map[string][]string) {
    wordlist, err := os.Open(wordlistPath)
    if err != nil {
        log.Fatal(err)
    }
    defer wordlist.Close()

    var subdomainMap = make(map[string][]string)

    scanner := bufio.NewScanner(wordlist)
    for scanner.Scan() {
        subdomain := scanner.Text() + "." + target;
        subdomain = strings.Replace(subdomain, " ", "", -1);

        ips, err := ResolveDNS(subdomain)
        if err == nil {
            subdomainMap[subdomain] = ips
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return subdomainMap
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
