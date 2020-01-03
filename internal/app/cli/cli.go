package cli

import (
	"fmt"
	"log"
	"os"

	"text/tabwriter"

	"github.com/leobeosab/sharingan/pkg/dns"
	"github.com/leobeosab/sharingan/pkg/nmap"
	"github.com/urfave/cli"
)

func SetupCLI() {

	var dnsWordlist string
	var target string

	sharingan := cli.NewApp()
	sharingan.Name = "Sharingan"
	sharingan.Usage = "Wrapper and analyzer for offensive security recon tools"

	sharingan.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "dns-wordlist",
			Value:       "",
			Usage:       "Wordlist for DNS bruteforcing",
			Destination: &dnsWordlist,
		},
		&cli.StringFlag{
			Name:        "target",
			Value:       "",
			Usage:       "Target domain",
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

	if target == "" {
		log.Fatal("Target needs to be defined")
	}

	if wordlistPath == "" {
		log.Fatal("DNS Wordlist needs to be defined")
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	hosts := dns.DNSBruteForce(target, wordlistPath)

	nmap.FilterHosts(&hosts)

	fmt.Fprintf(w, "\n%s\t%s\t", "Host Address ", "| Subdomain List")
	fmt.Fprintf(w, "\n%s\t%s\t", "----------------- ", "| -------")
	for _, h := range hosts {
		fmt.Fprintf(w, "\n%s \t| %v\t", h.IP, h.Subdomains)
	}

	w.Flush()
}
