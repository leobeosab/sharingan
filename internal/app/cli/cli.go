package cli

import (
	"log"
    "os"

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
        dns.RunDNSRecon(target, dnsWordlist)
		return nil
	}

	err := sharingan.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
