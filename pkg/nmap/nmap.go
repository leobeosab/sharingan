package nmap

import (
	"context"
	"log"
	"time"

	"github.com/Ullaakut/nmap"
	"github.com/leobeosab/sharingan/internal/models"
)

func Scan(target string) []models.Port {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithContext(ctx),
		nmap.WithFastMode(),
	)

	if err != nil {
		log.Fatalf("Unable to create nmap scanner: %v", err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("Unable to run nmap scan: %v", err)

		if len(result.NmapErrors) > 0 {
			for e := range result.NmapErrors {
				log.Println(e)
			}
		}
	}

	ports := make([]models.Port, 0)

	if len(result.Hosts) == 0 {
		return ports
	}
	// No support for multiple hosts at once yet
	for _, np := range result.Hosts[0].Ports {
		p := models.Port{
			ID:          np.ID,
			Protocol:    np.Protocol,
			ServiceName: np.Service.Name,
		}

		ports = append(ports, p)
	}

	return ports
}
