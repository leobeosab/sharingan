package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/leobeosab/sharingan/internal/helpers"
	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/dns"
	"github.com/leobeosab/sharingan/pkg/storage"
)

// Runs DNS recon and adds found hosts to program
func RunDNSRecon(settings *models.ScanSettings) {
	if settings.Target == "" {
		log.Fatal("Target needs to be defined")
	}
	if settings.DNSWordlist == "" {
		log.Fatalf("No program found - DNS Wordlist needs to be defined")
	}

	var p models.Program

	r := storage.RetrieveProgram(settings.Store, settings.Target)
	if len(r) != 0 {
		p = r[0]
	} else {
		p = models.Program{
			ProgramName: settings.Target,
		}
	}

	subs := dns.DNSBruteForce(settings.RootDomain, settings.DNSWordlist, settings.Threads)
	// Pesky progressbars not ending their lines
	fmt.Printf("\n")

	if settings.Rescan {
		AddSubsToProgram(&p, &subs)
	} else {
		ReplaceSubsInProgram(&p, &subs)
	}

	storage.UpdateOrCreateProgram(settings.Store, &p)
}

// Reads from STD in line by line adding hosts to program
func AddSubsFromInput(settings *models.ScanSettings) {

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
			Hosts:       make(map[string]models.Host),
		}
	} else {
		p = r[0]
	}

	// If there is no input in stdin, print an error
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
	subdomains = helpers.RemoveDuplicatesInSlice(subdomains)

	cliOut := fmt.Sprintf("Added %v subdomains to %s \n", len(subdomains), settings.Target)

	if settings.ReplaceSubs {
		ReplaceSubsInProgram(&p, &subdomains)
		cliOut = fmt.Sprintf("Replacing subdomains for %s \n", settings.Target)
	} else {
		AddSubsToProgram(&p, &subdomains)
	}

	storage.UpdateOrCreateProgram(settings.Store, &p)

	log.Printf(cliOut)
}

// Remove slice of subs for another
func ReplaceSubsInProgram(p *models.Program, subs *[]string) {
	p.Hosts = make(map[string]models.Host)
	for _, s := range *subs {
		fmt.Println(s)
		p.Hosts[s] = models.Host{
			Subdomain: s,
		}
	}
}

// Add subs to a fiven program
func AddSubsToProgram(p *models.Program, subs *[]string) {
	for _, s := range *subs {
		if _, ok := p.Hosts[s]; ok {
			continue
		}

		fmt.Println(s)
		p.Hosts[s] = models.Host{
			Subdomain: s,
		}
	}
}
