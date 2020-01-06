package dns

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/leobeosab/sharingan/internal/helpers"
	"github.com/leobeosab/sharingan/internal/models"
	"github.com/schollz/progressbar/v2"
)

func DNSBruteForce(target string, wordlistPath string) []models.Host {
	// Read the DNS names from wordlist
	wordlist, err := os.Open(wordlistPath)
	if err != nil {
		log.Fatal(err)
	}
	defer wordlist.Close()

	// Output information to the users
	fmt.Printf("\n\nBeginnning DNS Brute Force\n")
	// Progress bar :: Needs refactoring
	lines := helpers.GetNumberOfLinesInFile(wordlist)
	progress := progressbar.NewOptions(lines, progressbar.OptionSetPredictTime(false))

	// map[host][]subdomains
	var hostmap = make(map[string][]string)

	mux := &sync.Mutex{}
	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for domain := range jobs {
				ip := ResolveDNS(domain)
				if ip != "Error" {
					mux.Lock()
					// if ip exists in hostmap do the following
					if _, ok := hostmap[ip]; ok {
						hostmap[ip] = append(hostmap[ip], domain)
					} else {
						hostmap[ip] = []string{domain}
					}
					mux.Unlock()
				}

				progress.Add(1)
			}
		}()
	}

	// Stream wordlist to resolve DNS
	wlstream := bufio.NewScanner(wordlist)
	for wlstream.Scan() {
		subdomain := wlstream.Text() + "." + target
		subdomain = strings.Replace(subdomain, " ", "", -1)
		jobs <- subdomain
	}

	close(jobs)
	wg.Wait()

	if err := wlstream.Err(); err != nil {
		log.Fatal(err)
	}

	return HostmapToHostSlice(hostmap)
}

// map[host][]subdomains
func HostmapToHostSlice(m map[string][]string) []models.Host {
	s := make([]models.Host, 0)

	for k, v := range m {
		h := &models.Host{
			IP:         k,
			Subdomains: v,
		}

		s = append(s, *h)
	}

	return s
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
