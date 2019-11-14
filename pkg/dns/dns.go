package dns

import (
	"context"
	"log"
	"os"
	"time"
    "bufio"
    "net"
    "strings"
)

func DNSBruteForce(target string, wordlistPath string) (map[string][]string) {
    wordlist, err := os.Open(wordlistPath)
    if err != nil {
        log.Fatal(err)
    }
    defer wordlist.Close()

    var subdomainMap = make(map[string][]string)

    scanner := bufio.NewScanner(wordlist)
    for scanner.Scan() {
        subdomain := scanner.Text() + "." + target;
        subdomain = strings.Replace(subdomain, " ", "", -1);

        ips, err := ResolveDNS(subdomain)
        if err == nil {
            subdomainMap[subdomain] = ips
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return subdomainMap
}

func ResolveDNS(subdomain string) ([]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    return net.DefaultResolver.LookupHost(ctx, subdomain)
}
