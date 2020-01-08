package cli

import (
	"fmt"
	"strings"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/nmap"
	"github.com/leobeosab/sharingan/pkg/storage"
	"github.com/manifoldco/promptui"
)

func RunNmapScan(s *models.ScanSettings) {

	// This feels gross, find a good way to return a single entry with bolthold
	exists, results := storage.ScanEntryExists(s.Store, s.Target)

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
