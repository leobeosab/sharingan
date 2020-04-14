package dns

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/leobeosab/sharingan/internal/helpers"
	"github.com/schollz/progressbar/v2"
)

func DNSBruteForce(rd string, wordlistPath string, threads int) []string {
	// Read the DNS names from wordlist
	wordlist, err := os.Open(wordlistPath)
	if err != nil {
		log.Fatal(err)
	}
	defer wordlist.Close()

	// Output information to the users
	log.Printf("Beginnning DNS Brute Force\n")
	// TODO: Create helper methods for progress bars and prompts
	// Progress bar :: Needs refactoring
	lines := helpers.GetNumberOfLinesInFile(wordlist)
	progress := progressbar.NewOptions(lines, progressbar.OptionSetPredictTime(false))

	// Process jobs with async
	jobs := make(chan string, lines)
	subdomains := make(chan string, lines)
	var wg sync.WaitGroup

	// Create goroutines/"threads" and process incoming jobs
	for i := 0; i < threads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for domain := range jobs {
				ip := ResolveDNS(domain)
				if ip != "Error" {
					subdomains <- domain
				}

				progress.Add(1)
			}
		}()
	}

	// Add possible subdomains to jobs list
	wlstream := bufio.NewScanner(wordlist)
	jobs <- rd
	for wlstream.Scan() {
		subdomain := wlstream.Text() + "." + rd
		subdomain = strings.Replace(subdomain, " ", "", -1)
		jobs <- subdomain
	}

	if err := wlstream.Err(); err != nil {
		log.Fatal(err)
	}

	// Cleanup
	close(jobs)
	wg.Wait() // Wait for all threads to complete
	close(subdomains)

	// Get all subdomains into a slice and remove dupes
	ss := make([]string, 0)
	for s := range subdomains {
		ss = append(ss, s)
	}
	ss = helpers.RemoveDuplicatesInSlice(ss)

	return ss
}

// We only return the first ip because we aren't interested in redundant servers (ie AWS ELB instances
func ResolveDNS(subdomain string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ip := "Error"

	ips, err := net.DefaultResolver.LookupHost(ctx, subdomain)
	if err == nil {
		ip = ips[0]
	}

	return ip
}
