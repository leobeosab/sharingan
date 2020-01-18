package cli

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/nmap"
	"github.com/leobeosab/sharingan/pkg/storage"
	"github.com/manifoldco/promptui"
)

func RunNmapScan(s *models.ScanSettings) {
	exists, results := storage.ProgramEntryExists(s.Store, s.Target)

	if exists {
		p := results[0]
		log.Printf("Starting Nmap scan of %v domains... this may take some time", len(p.Hosts))

		hosts := make(chan models.Host, len(p.Hosts))
		results := make(chan models.Host, len(p.Hosts))
		var wg sync.WaitGroup

		for t := 0; t < s.Threads; t++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for h := range hosts {
					h.Ports = nmap.Scan(h.IP)
					log.Printf("\n%v : %v\n", h.IP, h.Ports)
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

		slice := make([]models.Host, 0)
		for h := range results {
			slice = append(slice, h)
		}

		p.Hosts = slice

		storage.UpdateProgram(s.Store, &p)

	} else {
		log.Printf("Error: no program called %s found \n", s.Target)
	}

}

func RunNmapScanInteractive(s *models.ScanSettings) {

	// This feels gross, find a good way to return a single entry with bolthold
	exists, results := storage.ProgramEntryExists(s.Store, s.Target)

	if exists {
		result := results[0]
		options := make([]string, 0)

		for _, h := range result.Hosts {
			option := h.IP + " - [" + strings.Join(h.Subdomains, ",") + "]"
			options = append(options, option)
		}

		prompt := promptui.Select{
			Label: "Select host to scan",
			Items: options,
		}

		_, selection, err := prompt.Run()

		if err != nil {
			fmt.Println("Error based on input")
		}

		d := strings.Split(selection, " - ")[0]
		fmt.Printf("Scanning %s with nmap...\n\n", d)
		nmap.Scan(d)

	} else {
		fmt.Printf("No scans found for %s", s.Target)

		prompt := promptui.Prompt{
			Label:     "Do you want to run the scan on: " + s.Target,
			IsConfirm: true,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}

		if result == "y" {
			nmap.Scan(s.Target)
		} else {
			fmt.Println("You got it champ, see ya")
		}
	}
}
