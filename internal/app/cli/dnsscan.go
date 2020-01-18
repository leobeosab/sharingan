package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/leobeosab/sharingan/internal/helpers"
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

		// We will need a way to update this instead of just overwriting it
		p.Hosts, p.Subdomains = dns.DNSBruteForce(settings.RootDomain, settings.DNSWordlist)

		if !settings.SkipProbe {
			nmap.FilterHosts(&p.Hosts)
		}

		if settings.Rescan {
			storage.UpdateProgram(settings.Store, &p)
		} else {
			storage.SaveProgram(settings.Store, &p)
		}
	} else {
		p = r[0]
		log.Println("Found previous scan")
	}

	// Get a few lines of space before final result
	log.Printf("\n")
	// Pretty print the results
	w := new(tabwriter.Writer)
	w.Init(os.Stderr, 8, 8, 0, '\t', 0)

	fmt.Fprintf(w, "\n%s\t%s\t", "Host Address ", "| Subdomain List")
	fmt.Fprintf(w, "\n%s\t%s\t", "----------------- ", "| -------")
	for _, h := range p.Hosts {
		fmt.Fprintf(w, "\n%s \t| %v\t", h.IP, h.Subdomains)
	}

	w.Flush()
}

func AddSubsToProgram(settings *models.ScanSettings) {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	var p models.Program
	r := storage.RetrieveProgram(settings.Store, settings.Target)
	if len(r) == 0 {
		fmt.Println(settings.Target + " not found in store, creating new entry")
		p = models.Program{
			ProgramName: settings.Target,
			Hosts:       []models.Host{},
		}
	} else {
		p = r[0]
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		log.Println("DNS addsubs is intended to work with pipies.")
		log.Println("Usage: cat subs | sharingancli --target program dns addsubs")
		return
	}

	reader := bufio.NewScanner(os.Stdin)
	var subdomains []string

	for reader.Scan() {
		input := reader.Text()
		subdomains = append(subdomains, input)
	}

	cliOut := fmt.Sprintf("Added %v subdomains to %s \n", len(subdomains), settings.Target)

	if !settings.ReplaceSubs {
		p.Subdomains = subdomains
	} else {
		p.Subdomains = append(p.Subdomains, subdomains...)
		cliOut = fmt.Sprintf("Replacing subdomains for %s \n", settings.Target)
	}
	p.Subdomains = helpers.RemoveDuplicatesInSlice(p.Subdomains)

	storage.UpdateOrCreateProgram(settings.Store, &p)

	log.Printf(cliOut)
}
