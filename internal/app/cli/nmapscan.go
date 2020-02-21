package cli

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/leobeosab/sharingan/internal/helpers"
	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/nmap"
	"github.com/leobeosab/sharingan/pkg/storage"
	"github.com/manifoldco/promptui"
)

func RunNmapScan(s *models.ScanSettings) {
	exists, p := storage.RetrieveOrCreateProgram(s.Store, s.Target)

	if exists {
		log.Printf("Starting Nmap scan of %v hosts... this may take some time\n", len(p.Hosts))

		hosts := make(chan models.Host, len(p.Hosts))
		results := make(chan models.Host, len(p.Hosts))
		var wg sync.WaitGroup

		if len(p.Hosts) > 10 && !s.NoPrompt {
			prompt := promptui.Prompt{
				Label:     "Warning port scans can be loud and you are scanning with > 10 hosts. Do you want to continue?",
				IsConfirm: true,
			}

			result, err := prompt.Run()

			if err != nil {
				log.Printf("Prompt failed %v\n", err)
				return
			}

			if result != "y" {
				log.Printf("Exiting... \n")
				return
			}
		}

		// This can be optimized to check if it has already scanned the same host...
		// under a different subdomain
		for t := 0; t < s.Threads; t++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for h := range hosts {
					h.Ports = nmap.Scan(h.Subdomain)
					helpers.PrintNmapScan(h)
					results <- h
				}
			}()
		}

		for _, h := range p.Hosts {
			hosts <- h
		}

		close(hosts)
		wg.Wait()
		close(results)

		for h := range results {
			p.Hosts[h.Subdomain] = h
		}

		storage.UpdateProgram(s.Store, &p)

	} else {
		log.Printf("Error: no program called %s found \n", s.Target)
	}

}

func RunNmapScanInteractive(s *models.ScanSettings) {
	exists, result := storage.RetrieveOrCreateProgram(s.Store, s.Target)

	if exists {
		options := make([]string, 0)

		for _, h := range result.Hosts {
			option := h.Subdomain
			options = append(options, option)
		}

		prompt := promptui.Select{
			Label: "Select host to scan",
			Items: options,
		}

		_, selection, err := prompt.Run()

		if err != nil {
			log.Printf("Prompt error\n")
			return
		}

		d := strings.Split(selection, " - ")[0]
		log.Printf("Scanning %s with nmap...\n\n", d)
		ports := nmap.Scan(d)

		if len(ports) == 0 {
			log.Printf("No ports open \n")
			return
		}

		host := result.Hosts[d]
		host.Ports = ports
		result.Hosts[d] = host

		helpers.PrintNmapScan(host)

		storage.UpdateProgram(s.Store, &result)
	} else {
		fmt.Printf("No scans found for %s", s.Target)
	}
}
