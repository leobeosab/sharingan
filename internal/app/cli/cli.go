package cli

import (
	"log"
    "os"
    "fmt"

    "text/tabwriter"

	"github.com/urfave/cli"
    "github.com/leobeosab/sharingan/pkg/dns"
)

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

    validSubdomains := dns.DNSBruteForce(target, wordlistPath)

    fmt.Fprintf(w, "\n%s\t%s\t", "Subdomain Address ", "| IP List")
    fmt.Fprintf(w, "\n%s\t%s\t", "----------------- ", "| -------")
    for subdomain, ips := range validSubdomains {
        fmt.Fprintf(w, "\n%s \t| %v\t", subdomain, ips)
    }

    w.Flush()
}
