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

	if settings.DNSWordlist == "" {
		log.Fatal("DNS Wordlist needs to be defined")
	}

	var s models.ScanResults

	r := storage.RetrieveScanResults(settings.Store, settings.Target)
	if len(r) == 0 || settings.Rescan {
		s = models.ScanResults{
			RootDomain: settings.Target,
		}

		s.Hosts = dns.DNSBruteForce(settings.Target, settings.DNSWordlist)

		if !settings.SkipProbe {
			nmap.FilterHosts(&s.Hosts)
		}

		if settings.Rescan {
			storage.UpdateScan(settings.Store, &s)
		} else {
			storage.SaveScan(settings.Store, &s)
		}
	} else {
		s = r[0]
		fmt.Println("Found previous scan")
	}

	// Get a few lines of space before final result
	fmt.Printf("\n\n")
	// Pretty print the results
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	fmt.Fprintf(w, "\n%s\t%s\t", "Host Address ", "| Subdomain List")
	fmt.Fprintf(w, "\n%s\t%s\t", "----------------- ", "| -------")
	for _, h := range s.Hosts {
		fmt.Fprintf(w, "\n%s \t| %v\t", h.IP, h.Subdomains)
	}

	w.Flush()
}
