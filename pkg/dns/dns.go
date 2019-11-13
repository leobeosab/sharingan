package dns

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
    "bufio"
    "net"
    "strings"

    "text/tabwriter"
)

func RunDNSRecon(target string, wordlistPath string) {

    if (target == "") {
        log.Fatal("Target needs to be defined")
    }

    if (wordlistPath == "") {
        log.Fatal("DNS Wordlist needs to be defined")
    }

    w := new(tabwriter.Writer)
    w.Init(os.Stdout, 8, 8, 0, '\t', 0)

    validSubdomains := DNSBruteForce(target, wordlistPath)

    fmt.Fprintf(w, "\n%s\t%s\t", "Subdomain Address ", "| IP List")
    fmt.Fprintf(w, "\n%s\t%s\t", "----------------- ", "| -------")
    for subdomain, ips := range validSubdomains {
        fmt.Fprintf(w, "\n%s \t| %v\t", subdomain, ips)
    }

    w.Flush()
}

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
