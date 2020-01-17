package cli

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/dns"
	"github.com/leobeosab/sharingan/pkg/nmap"
	"github.com/leobeosab/sharingan/pkg/storage"
)

func RunDNSRecon(settings *models.ScanSettings) {
	if settings.Target == "" {
		log.Fatal("Target needs to be defined")
	}

	var p models.Program

	r := storage.RetrieveProgram(settings.Store, settings.Target)
	if len(r) == 0 || settings.Rescan {
		// Check for DNS wordlist
		if settings.DNSWordlist == "" {
			log.Fatalf("No program found - DNS Wordlist needs to be defined")
		}

		p = models.Program{
			ProgramName: settings.Target,
		}

		p.Hosts = dns.DNSBruteForce(settings.RootDomain, settings.DNSWordlist)

		if !settings.SkipProbe {
			nmap.FilterHosts(&p.Hosts)
		}

		if settings.Rescan {
			storage.UpdateScan(settings.Store, &p)
		} else {
			storage.SaveProgram(settings.Store, &p)
		}
	} else {
		p = r[0]
		fmt.Println("Found previous scan")
	}

	// Get a few lines of space before final result
	fmt.Printf("\n")
	// Pretty print the results
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	fmt.Fprintf(w, "\n%s\t%s\t", "Host Address ", "| Subdomain List")
	fmt.Fprintf(w, "\n%s\t%s\t", "----------------- ", "| -------")
	for _, h := range p.Hosts {
		fmt.Fprintf(w, "\n%s \t| %v\t", h.IP, h.Subdomains)
	}

	w.Flush()
}
