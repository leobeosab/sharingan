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
	sharingan.Version = "0.1.0"

	sharingan.Flags = GetGlobalFlags(settings)
	sharingan.Commands = []*cli.Command{
		{
			Name:  "dns",
			Usage: "Perform a DNS scan : sharingancli --target target.com dns --dns-wordlist ./path/to/list",
			Flags: GetDNSFlags(settings),
			Action: func(c *cli.Context) error {
				RunDNSRecon(settings)
				return nil
			},
		},
		{
			Name:  "scan",
			Usage: "Perform a service scan using nmap -sV : sharingancli --target target.com scan",
			Flags: GetDNSFlags(settings),
			Action: func(c *cli.Context) error {
				RunNmapScan(settings)
				return nil
			},
		},
	}

	sharingan.Action = func(c *cli.Context) error {
		return nil
	}

	err := sharingan.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
