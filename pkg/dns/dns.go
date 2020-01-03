package dns

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/leobeosab/sharingan/internal/models"
)

func DNSBruteForce(target string, wordlistPath string) []models.Host {
	// Read the DNS names from wordlist
	wordlist, err := os.Open(wordlistPath)
	if err != nil {
		log.Fatal(err)
	}
	defer wordlist.Close()

	// map[host][]subdomains
	var hostmap = make(map[string][]string)

	// Stream wordlist to resolve DNS
	wlstream := bufio.NewScanner(wordlist)
	for wlstream.Scan() {
		subdomain := wlstream.Text() + "." + target
		subdomain = strings.Replace(subdomain, " ", "", -1)

		ip := ResolveDNS(subdomain)
		if ip != "Error" {
			// if ip exists in hostmap do the following
			if _, ok := hostmap[ip]; ok {
				hostmap[ip] = append(hostmap[ip], subdomain)
			} else {
				hostmap[ip] = []string{subdomain}
			}
		}
	}

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
