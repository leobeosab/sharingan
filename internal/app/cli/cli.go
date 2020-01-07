package cli

import (
	"log"
	"os"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/storage"
	"github.com/urfave/cli/v2"
)

func SetupCLI() {

	settings := &models.ScanSettings{}

	settings.Store = storage.OpenStore()
	defer settings.Store.Close()

	sharingan := cli.NewApp()
	sharingan.Name = "Sharingan"
	sharingan.Usage = "Wrapper and analyzer for offensive security recon tools"

	sharingan.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "dns-wordlist",
			Value:       "",
			Usage:       "Wordlist for DNS bruteforcing",
			Destination: &settings.DNSWordlist,
		},
		&cli.StringFlag{
			Name:        "target",
			Value:       "",
			Usage:       "Target domain",
			Destination: &settings.Target,
		},
		&cli.BoolFlag{
			Name:        "skip-probe",
			Usage:       "Skips host-up nmap scan",
			Destination: &settings.SkipProbe,
		},
		&cli.BoolFlag{
			Name:        "rescan",
			Usage:       "Scans domain regardless of the existance of previous results",
			Destination: &settings.Rescan,
		},
	}

	sharingan.Commands = []*cli.Command{
		{
			Name:  "dns",
			Usage: "Perform a DNS scan",
			Action: func(c *cli.Context) error {
				RunDNSRecon(settings)
				return nil
			},
		},
		{
			Name:  "scan",
			Usage: "Perform a service scan using nmap -sV",
			Action: func(c *cli.Context) error {
				RunNmapScan(settings)
				return nil
			},
		},
	}

	sharingan.Action = func(c *cli.Context) error {
		// RunDNSRecon(settings)
		return nil
	}

	err := sharingan.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
